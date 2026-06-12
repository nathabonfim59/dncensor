package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/nathabonfim59/dncensor/internal/backup"
	"github.com/nathabonfim59/dncensor/internal/dns"
	"github.com/nathabonfim59/dncensor/internal/provider"
	"github.com/nathabonfim59/dncensor/internal/stack"
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

	// Flavor menu
	flavorIdx int

	// Result
	result *dns.ApplyResult

	// Error
	err error

	// Window
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
				// providers + apply button
				maxIdx := len(m.providers) // apply button is last
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
	// Check if we're on the apply button (last item)
	if m.selectedIdx == len(m.providers) {
		// Apply
		if m.selectedProvider == nil {
			m.err = fmt.Errorf("no provider selected")
			return m, nil
		}
		m.state = stateConfirm
		return m, nil
	}

	p := m.providers[m.selectedIdx]
	m.selectedProvider = p

	if p.SupportsFlavors() && m.selectedIdx == 1 { // CloudFlare is index 1
		m.flavorIdx = 0
		m.selectedFlavor = nil
		m.state = stateFlavorMenu
		return m, nil
	}

	// For ISP and Google, just select (no flavor menu)
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

func (m Model) View() string {
	switch m.state {
	case stateMainMenu:
		return m.mainMenuView()
	case stateFlavorMenu:
		return m.flavorMenuView()
	case stateConfirm:
		return m.confirmView()
	case stateResult:
		return m.resultView()
	case stateExiting:
		return ""
	}
	return ""
}

func (m Model) mainMenuView() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render(fmt.Sprintf(" 🛡 %s — DNS Provider Switcher\n\n", "dncensor")))
	b.WriteString(LabelStyle.Render(fmt.Sprintf("Detected stack: %s\n\n", m.stack.Type())))

	// Provider list
	b.WriteString(LabelStyle.Render("Select DNS Provider:\n"))
	for i, p := range m.providers {
		cursor := "  "
		selected := " "
		if i == m.selectedIdx {
			cursor = "▸ "
		}
		if m.selectedProvider == p {
			selected = "●"
		} else {
			selected = "○"
		}

		// Grey out ISP if no backup
		entry := fmt.Sprintf("%s %s %s", cursor, selected, p.Name)
		if p.Type == provider.ProviderISP && m.backupManager.Exists(string(m.stack.Type())) == false {
			// Try to check if backup exists by scanning
			if !m.hasBackup() {
				b.WriteString(DimmedStyle.Render(fmt.Sprintf("%s %s %s (no backup)", cursor, "○", p.Name)))
				b.WriteString("\n")
				continue
			}
		}

		if i == m.selectedIdx {
			b.WriteString(SelectedStyle.Render(entry))
		} else {
			b.WriteString(NormalStyle.Render(entry))
		}
		b.WriteString("\n")
	}

	// Apply button
	b.WriteString("\n")
	applyText := "  [ Apply Configuration ]"
	if len(m.providers) == m.selectedIdx {
		b.WriteString(SelectedStyle.Render(applyText))
	} else {
		b.WriteString(NormalStyle.Render(applyText))
	}

	// DoH toggle
	b.WriteString("\n\n")
	dohLabel := "  Use DNS-over-HTTPS (DoH)"
	if m.useDOH {
		dohLabel = CheckboxChecked.Render("  ☑ Use DNS-over-HTTPS (DoH)")
	} else {
		dohLabel = CheckboxStyle.Render("  ☐ Use DNS-over-HTTPS (DoH)")
	}
	b.WriteString(dohLabel)
	b.WriteString("\n")

	// Show selected provider info
	if m.selectedProvider != nil {
		b.WriteString("\n")
		b.WriteString(DimmedStyle.Render(fmt.Sprintf("  Selected: %s", m.selectedProvider.Name)))
		if m.selectedFlavor != nil {
			b.WriteString(DimmedStyle.Render(fmt.Sprintf(" > %s", m.selectedFlavor.Display)))
		}
		b.WriteString("\n")
	}

	if m.err != nil {
		b.WriteString("\n")
		b.WriteString(ErrorStyle.Render(fmt.Sprintf("  Error: %s", m.err)))
		m.err = nil
	}

	b.WriteString("\n")
	b.WriteString(HintStyle.Render("  ↑/↓ navigate • enter select • tab toggle DoH • a apply • q quit"))

	return BorderStyle.Render(b.String())
}

func (m Model) hasBackup() bool {
	_, err := m.backupManager.Latest(string(m.stack.Type()))
	return err == nil
}

func (m Model) flavorMenuView() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render(fmt.Sprintf(" ☁ %s — Choose Flavor\n\n", "CloudFlare")))

	b.WriteString(LabelStyle.Render("Select flavor:\n"))
	for i, f := range m.selectedProvider.Flavors {
		cursor := "  "
		if i == m.flavorIdx {
			cursor = "▸ "
		}
		entry := fmt.Sprintf("%s %s (%s)", cursor, f.Display, f.PrimaryDNS)

		if i == m.flavorIdx {
			b.WriteString(SelectedStyle.Render(entry))
		} else {
			b.WriteString(NormalStyle.Render(entry))
		}
		b.WriteString("\n")
		b.WriteString(DimmedStyle.Render(fmt.Sprintf("    %s\n", f.Description)))
	}

	b.WriteString("\n")
	b.WriteString(HintStyle.Render("  ↑/↓ navigate • enter select • esc back"))

	return BorderStyle.Render(b.String())
}

func (m Model) confirmView() string {
	var b strings.Builder

	b.WriteString(TitleStyle.Render(" ⚡ Confirm DNS Change\n\n"))

	b.WriteString(fmt.Sprintf("%s %s\n", LabelStyle.Render("Provider:"), ValueStyle.Render(m.selectedProvider.Name)))
	if m.selectedFlavor != nil {
		b.WriteString(fmt.Sprintf("%s %s\n", LabelStyle.Render("Flavor:"), ValueStyle.Render(m.selectedFlavor.Display)))
	}
	b.WriteString(fmt.Sprintf("%s %v\n", LabelStyle.Render("DoH:"), ValueStyle.Render(fmt.Sprintf("%v", m.useDOH))))

	primary, secondary, _, _ := m.selectedProvider.Resolve("", m.useDOH)
	if m.selectedFlavor != nil {
		primary, secondary, _, _ = m.selectedProvider.Resolve(m.selectedFlavor.FlavorName, m.useDOH)
	}
	b.WriteString(fmt.Sprintf("%s %s / %s\n", LabelStyle.Render("DNS Servers:"), ValueStyle.Render(primary), ValueStyle.Render(secondary)))

	b.WriteString("\n")
	b.WriteString(LabelStyle.Render("Press enter to apply, esc to cancel"))

	return BorderStyle.Render(b.String())
}

func (m Model) resultView() string {
	var b strings.Builder

	if m.result.Success {
		b.WriteString(SuccessStyle.Render(" ✓ DNS Configuration Applied\n\n"))
		b.WriteString(ValueStyle.Render(m.result.Message))
	} else {
		b.WriteString(ErrorStyle.Render(" ✗ Failed to Apply DNS\n\n"))
		b.WriteString(ValueStyle.Render(m.result.Message))
	}

	b.WriteString("\n\n")
	b.WriteString(HintStyle.Render("  enter • Back to menu   q • Quit"))

	return BorderStyle.Render(b.String())
}
