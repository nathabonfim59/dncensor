package tui

import (
	"fmt"

	"github.com/charmbracelet/bubbletea"
	"github.com/nathabonfim59/dncensor/internal/backup"
	"github.com/nathabonfim59/dncensor/internal/dns"
	"github.com/nathabonfim59/dncensor/internal/provider"
	"github.com/nathabonfim59/dncensor/internal/stack"

	"charm.land/lipgloss/v2"
)

type state int

const (
	stateMainMenu state = iota
	stateFlavorMenu
	stateConfirm
	stateResult
	stateExiting
)

type Model struct {
	state            state
	stack            stack.Stack
	backupDir        string
	backupManager    *backup.BackupManager

	providers        []*provider.DNSProvider
	selectedIdx      int
	selectedProvider *provider.DNSProvider
	selectedFlavor   *provider.DNSFlavor
	useDOH           bool

	flavorIdx int

	result *dns.ApplyResult

	err error

	width  int
	height int
}

func NewModel(s stack.Stack, backupDir string) Model {
	return Model{
		state:         stateMainMenu,
		stack:         s,
		backupDir:     backupDir,
		backupManager: backup.New(backupDir),
		providers:     provider.AllProviders(),
		selectedIdx:   0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

// ── Update ──────────────────────────────────────────────────────────────

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit

		case "esc":
			switch m.state {
			case stateFlavorMenu:
				m.state = stateMainMenu
			case stateConfirm:
				m.state = stateMainMenu
			case stateResult:
				m.state = stateMainMenu
			}
			return m, nil

		case "up", "k":
			switch m.state {
			case stateMainMenu:
				if m.selectedIdx > 0 {
					m.selectedIdx--
				}
			case stateFlavorMenu:
				if m.flavorIdx > 0 {
					m.flavorIdx--
				}
			}
			return m, nil

		case "down", "j":
			switch m.state {
			case stateMainMenu:
				maxIdx := len(m.providers)
				if m.selectedIdx < maxIdx {
					m.selectedIdx++
				}
			case stateFlavorMenu:
				if m.flavorIdx < len(m.selectedProvider.Flavors)-1 {
					m.flavorIdx++
				}
			}
			return m, nil

		case "enter", " ":
			switch m.state {
			case stateMainMenu:
				return m.handleMainMenuEnter()
			case stateFlavorMenu:
				m.selectedFlavor = &m.selectedProvider.Flavors[m.flavorIdx]
				m.state = stateMainMenu
				return m, nil
			case stateConfirm:
				return m.applyDNS()
			case stateResult:
				m.state = stateMainMenu
				m.result = nil
				m.selectedProvider = nil
				m.selectedFlavor = nil
				m.selectedIdx = 0
				return m, nil
			}
			return m, nil

		case "a", "A":
			if m.state == stateMainMenu {
				return m.handleMainMenuEnter()
			}
			return m, nil

		case "tab":
			if m.state == stateMainMenu {
				m.useDOH = !m.useDOH
			}
			return m, nil

		case "left", "b":
			if m.state == stateFlavorMenu || m.state == stateConfirm {
				m.state = stateMainMenu
			}
			return m, nil
		}
	}

	return m, nil
}

func (m Model) handleMainMenuEnter() (tea.Model, tea.Cmd) {
	if m.selectedIdx == len(m.providers) {
		if m.selectedProvider == nil {
			m.err = fmt.Errorf("no provider selected")
			return m, nil
		}
		m.state = stateConfirm
		return m, nil
	}

	p := m.providers[m.selectedIdx]
	m.selectedProvider = p

	if p.SupportsFlavors() {
		m.flavorIdx = 0
		m.selectedFlavor = nil
		m.state = stateFlavorMenu
		return m, nil
	}

	m.selectedFlavor = nil
	return m, nil
}

func (m Model) applyDNS() (tea.Model, tea.Cmd) {
	cfg := dns.ApplyConfig{
		Provider:   m.selectedProvider,
		FlavorName: "",
		UseDOH:     m.useDOH,
	}

	if m.selectedFlavor != nil {
		cfg.FlavorName = m.selectedFlavor.FlavorName
	}

	result := dns.Apply(m.stack, cfg, m.backupDir)
	m.result = &result
	m.state = stateResult
	return m, nil
}

func (m *Model) hasBackup() bool {
	_, err := m.backupManager.Latest(string(m.stack.Type()))
	return err == nil
}

// ── View ────────────────────────────────────────────────────────────────

func (m Model) View() string {
	w, h := m.width, m.height
	if w == 0 {
		w = 80
	}
	if h == 0 {
		h = 24
	}

	switch m.state {
	case stateMainMenu:
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, m.mainMenuView())
	case stateFlavorMenu:
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, m.flavorMenuView())
	case stateConfirm:
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, m.confirmView())
	case stateResult:
		return lipgloss.Place(w, h, lipgloss.Center, lipgloss.Center, m.resultView())
	case stateExiting:
		return ""
	}
	return ""
}

// ── Main menu ───────────────────────────────────────────────────────────

func (m Model) mainMenuView() string {
	sections := []string{
		m.renderHeader(),
		"",
		m.renderProviderList(),
		"",
		m.renderApplyButton(),
		"",
		m.renderDohToggle(),
		"",
		m.renderSelectionInfo(),
		m.renderError(),
		HintStyle.Render("↑/↓ navigate · enter select · tab toggle DoH · a apply · q quit"),
	}

	body := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return CardStyle.Render(body)
}

func (m Model) renderHeader() string {
	title := TitleStyle.Render("🛡 dncensor — DNS Provider Switcher")
	subtitle := DimmedStyle.Render(fmt.Sprintf("Detected stack: %s", m.stack.Type()))
	return lipgloss.JoinVertical(lipgloss.Left, title, subtitle)
}

func (m Model) renderProviderList() string {
	var items []string
	for i, p := range m.providers {
		cursor, radio := " ", "○"
		if i == m.selectedIdx {
			cursor = "▸"
		}
		if m.selectedProvider == p {
			radio = "●"
		}
		text := fmt.Sprintf("%s %s %s", cursor, radio, p.Name)

		if p.Type == provider.ProviderISP && !m.hasBackup() {
			items = append(items, ItemDimmedStyle.Render(text+" (no backup)"))
			continue
		}
		if i == m.selectedIdx {
			items = append(items, ItemSelectedStyle.Render(text))
		} else {
			items = append(items, ItemNormalStyle.Render(text))
		}
	}
	list := lipgloss.JoinVertical(lipgloss.Left, items...)
	return lipgloss.JoinVertical(lipgloss.Left,
		HeadingStyle.Render("Select DNS Provider:"),
		list,
	)
}

func (m Model) renderApplyButton() string {
	if m.selectedIdx == len(m.providers) {
		return lipgloss.PlaceHorizontal(26, lipgloss.Center,
			ApplyBtnStyle.Render("Apply Configuration"))
	}
	return lipgloss.PlaceHorizontal(26, lipgloss.Center,
		ApplyBtnDimmedStyle.Render("Apply Configuration"))
}

func (m Model) renderDohToggle() string {
	text := "Use DNS-over-HTTPS (DoH)"
	if m.useDOH {
		return CheckOnStyle.Render("☑ " + text)
	}
	return CheckOffStyle.Render("☐ " + text)
}

func (m Model) renderSelectionInfo() string {
	if m.selectedProvider == nil {
		return ""
	}
	info := fmt.Sprintf("Selected: %s", m.selectedProvider.Name)
	if m.selectedFlavor != nil {
		info += fmt.Sprintf(" > %s", m.selectedFlavor.Display)
	}
	return DimmedStyle.Render(info)
}

func (m Model) renderError() string {
	if m.err == nil {
		return ""
	}
	msg := fmt.Sprintf("Error: %s", m.err)
	m.err = nil
	return ErrorStyle.Render(msg)
}

// ── Flavor menu ─────────────────────────────────────────────────────────

func (m Model) flavorMenuView() string {
	sections := []string{
		TitleStyle.Render("☁ CloudFlare — Choose Flavor"),
		"",
		HeadingStyle.Render("Select flavor:"),
		m.renderFlavorList(),
		"",
		HintStyle.Render("↑/↓ navigate · enter select · esc back"),
	}

	body := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return CardStyle.Render(body)
}

func (m Model) renderFlavorList() string {
	var items []string
	for i, f := range m.selectedProvider.Flavors {
		prefix := "  "
		if i == m.flavorIdx {
			prefix = "▸ "
		}
		nameLine := fmt.Sprintf("%s%s (%s)", prefix, f.Display, f.PrimaryDNS)
		descLine := fmt.Sprintf("  %s", f.Description)

		nameStyle := ItemNormalStyle
		if i == m.flavorIdx {
			nameStyle = ItemSelectedStyle
		}
		entry := lipgloss.JoinVertical(lipgloss.Left,
			nameStyle.Render(nameLine),
			DimmedStyle.Render(descLine),
		)
		items = append(items, entry)
	}
	return lipgloss.JoinVertical(lipgloss.Left, items...)
}

// ── Confirm ─────────────────────────────────────────────────────────────

func (m Model) confirmView() string {
	primary, secondary, _, _ := m.selectedProvider.Resolve("", m.useDOH)
	if m.selectedFlavor != nil {
		primary, secondary, _, _ = m.selectedProvider.Resolve(m.selectedFlavor.FlavorName, m.useDOH)
	}

	rows := []string{
		m.renderConfirmRow("Provider:", m.selectedProvider.Name),
	}
	if m.selectedFlavor != nil {
		rows = append(rows, m.renderConfirmRow("Flavor:", m.selectedFlavor.Display))
	}
	rows = append(rows,
		m.renderConfirmRow("DoH:", fmt.Sprintf("%v", m.useDOH)),
		m.renderConfirmRow("DNS Servers:", fmt.Sprintf("%s / %s", primary, secondary)),
	)

	sections := []string{
		TitleStyle.Render("⚡ Confirm DNS Change"),
		"",
		lipgloss.JoinVertical(lipgloss.Left, rows...),
		"",
		HeadingStyle.Render("Press enter to apply, esc to cancel"),
	}

	body := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return CardStyle.Render(body)
}

func (m Model) renderConfirmRow(label, value string) string {
	return lipgloss.JoinHorizontal(lipgloss.Top,
		RowLabelStyle.Render(label),
		"  ",
		ValueStyle.Render(value),
	)
}

// ── Result ──────────────────────────────────────────────────────────────

func (m Model) resultView() string {
	var title string
	if m.result.Success {
		title = SuccessStyle.Render("✓ DNS Configuration Applied")
	} else {
		title = ErrorStyle.Render("✗ Failed to Apply DNS")
	}

	sections := []string{
		title,
		"",
		ValueStyle.Render(m.result.Message),
		"",
		HintStyle.Render("enter · Back to menu   q · Quit"),
	}

	body := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return CardStyle.Render(body)
}
