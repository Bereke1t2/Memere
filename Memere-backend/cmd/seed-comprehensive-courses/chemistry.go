package main

func buildChemistryCourse() CourseData {
	return CourseData{
		Title:            "Grade 12 Chemistry: Ethiopian National University Entrance Prep",
		Slug:             "grade-12-chemistry-euee-prep",
		Subject:          "Chemistry",
		Grade:            12,
		Description:      "Comprehensive Grade 11 & 12 Chemistry course focused on high performance in the Ethiopian University Entrance Examination. Master Atomic Structure & Periodicity, Chemical Kinetics & Le Chatelier Equilibrium, Acid-Base Buffer Calculations, Electrochemistry & Galvanic Cells, and Organic Reactions with structured quizzes and mock exams.",
		ShortDescription: "Grade 12 Chemistry EUEE prep covering atomic bonding, kinetics, acids/bases, electrochemistry, organic chemistry, quizzes, and exams.",
		ThumbnailUrl:     "https://images.unsplash.com/photo-1603126857599-f6e157fa2fe6?w=800&auto=format&fit=crop&q=80",
		Price:            0,
		IsFree:           true,
		Level:            "advanced",
		Sections: []SectionData{
			{
				Title:       "Unit 1: Atomic Structure, Quantum Numbers & Chemical Bonding",
				Description: "Quantum numbers (n, l, ml, ms), electron configuration, periodic trends, ionic vs covalent bonding, and VSEPR molecular geometry.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 1.1: Quantum Mechanical Model & Electronic Configuration",
						Type:            "video",
						DurationSeconds: 1250,
						IsFreePreview:   true,
						Content: `# Quantum Numbers and Orbitals
Four quantum numbers uniquely specify the state of an electron:
1. **Principal ($n$)**: Energy level ($n = 1, 2, 3...$).
2. **Azimuthal / Angular Momentum ($l$)**: Subshell shape ($l = 0 \to s, 1 \to p, 2 \to d, 3 \to f$).
3. **Magnetic ($m_l$)**: Orientation in space ($-l \le m_l \le +l$).
4. **Spin ($m_s$)**: $+\frac{1}{2}$ or $-\frac{1}{2}$.

### Rules for Filling Orbitals:
- **Aufbau Principle**: Lower energy orbitals fill first ($1s < 2s < 2p < 3s < 3p < 4s < 3d$).
- **Pauli Exclusion Principle**: No two electrons can have identical four quantum numbers.
- **Hund's Rule**: Degenerate orbitals are singly occupied with parallel spins before pairing.`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 1.2: VSEPR Theory & Intermolecular Forces",
						Type:            "note",
						DurationSeconds: 850,
						IsFreePreview:   true,
						Content: `# Molecular Geometry (VSEPR)
Electron pairs repel each other, adopting shapes that maximize distance:
- **Linear**: $180^\circ$ (e.g. $\text{BeCl}_2$, $\text{CO}_2$)
- **Trigonal Planar**: $120^\circ$ (e.g. $\text{BF}_3$)
- **Tetrahedral**: $109.5^\circ$ (e.g. $\text{CH}_4$)
- **Trigonal Pyramidal**: $107^\circ$ (e.g. $\text{NH}_3$)
- **Bent**: $104.5^\circ$ (e.g. $\text{H}_2\text{O}$)

### Intermolecular Forces:
Hydrogen bonding ($\text{H}$ bonded to $\text{F, O, N}$) $>$ Dipole-Dipole $>$ London Dispersion forces.`,
						PdfUrl: "sample.pdf",
					},
				},
			},
			{
				Title:       "Unit 2: Chemical Kinetics & Dynamic Equilibrium",
				Description: "Rate laws, reaction mechanisms, activation energy, collision theory, equilibrium constant (Kc, Kp), and Le Chatelier's Principle.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 2.1: Rate Laws & Arrhenius Activation Energy",
						Type:            "video",
						DurationSeconds: 1400,
						IsFreePreview:   false,
						Content: `# Chemical Kinetics
- **Rate Law**: $\text{Rate} = k [A]^m [B]^n$
- **Arrhenius Equation**:
  $$k = A e^{-E_a / RT} \implies \ln k = \ln A - \frac{E_a}{R}\left(\frac{1}{T}\right)$$
- Plotting $\ln k$ versus $1/T$ yields a straight line with slope $= -\frac{E_a}{R}$.`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 2.2: Equilibrium Constant & Le Chatelier's Principle",
						Type:            "note",
						DurationSeconds: 950,
						IsFreePreview:   false,
						Content: `# Dynamic Chemical Equilibrium
For reaction $aA + bB \rightleftharpoons cC + dD$:
$$K_c = \frac{[C]^c [D]^d}{[A]^a [B]^b}, \quad K_p = K_c (RT)^{\Delta n}$$

### Le Chatelier's Rules:
- **Adding Reactants**: Shifts equilibrium **Right** ($\to$).
- **Increasing Pressure**: Shifts to the side with **fewer moles of gas**.
- **Exothermic Reaction ($\Delta H < 0$) + Temperature Increase**: Shifts **Left** ($\leftarrow$) and decreases $K_c$.`,
						PdfUrl: "sample.pdf",
					},
				},
			},
			{
				Title:       "Unit 3: Acid-Base Equilibria, pH & Buffers",
				Description: "Brønsted-Lowry & Lewis definitions, strong/weak acids, Kw, Ka, Kb, Henderson-Hasselbalch equation, and titration curves.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 3.1: pH, pOH, and Weak Acid/Base Dissociation",
						Type:            "video",
						DurationSeconds: 1300,
						IsFreePreview:   false,
						Content: `# Acid-Base Dissociation
- **Ion-Product of Water**: $K_w = [H^+][OH^-] = 1.0 \times 10^{-14}$ at $25^\circ\text{C}$.
- $\text{pH} = -\log[H^+]$, $\quad \text{pOH} = -\log[OH^-]$, $\quad \text{pH} + \text{pOH} = 14$.
- **Weak Acid Dissociation**:
  $$K_a = \frac{[H^+][A^-]}{[HA]} \implies [H^+] = \sqrt{K_a \cdot C_a}$$`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 3.2: Buffer Solutions & Henderson-Hasselbalch Equation",
						Type:            "note",
						DurationSeconds: 900,
						IsFreePreview:   false,
						Content: `# Buffer Solutions
A buffer resists pH changes upon addition of small amounts of strong acid or base.

### Henderson-Hasselbalch Formula:
$$\text{pH} = \text{p}K_a + \log\left(\frac{[\text{Conjugate Base}]}{[\text{Weak Acid}]}\right)$$
Maximal buffering capacity occurs when $[\text{Base}] = [\text{Acid}]$, i.e., $\text{pH} = \text{p}K_a$.`,
						PdfUrl: "sample.pdf",
					},
				},
			},
			{
				Title:       "Unit 4: Electrochemistry & Galvanic Cells",
				Description: "Redox balancing, standard electrode potentials (E°), Nernst equation, electrolytic cells, and Faraday's laws of electrolysis.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 4.1: Galvanic Cells & Standard Cell Potentials",
						Type:            "video",
						DurationSeconds: 1350,
						IsFreePreview:   false,
						Content: `# Galvanic (Voltaic) Cells
- **Anode**: Site of **Oxidation** ($\text{Zn} \to \text{Zn}^{2+} + 2e^-$). Negative terminal.
- **Cathode**: Site of **Reduction** ($\text{Cu}^{2+} + 2e^- \to \text{Cu}$). Positive terminal.
- **Standard Cell EMF**:
  $$E^\circ_{\text{cell}} = E^\circ_{\text{cathode}} - E^\circ_{\text{anode}}$$
- Reaction is spontaneous when $E^\circ_{\text{cell}} > 0$ ($\Delta G^\circ = -n F E^\circ < 0$).`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 4.2: Faraday's Laws of Electrolysis & Nernst Equation",
						Type:            "note",
						DurationSeconds: 920,
						IsFreePreview:   false,
						Content: `# Electrolysis & Faraday's Law
$$m = \frac{M \cdot I \cdot t}{n \cdot F}$$
where $F = 96,500\text{ C/mol }e^-$, $I$ in Amperes, $t$ in seconds, $M$ = molar mass, $n$ = valence electrons.

### Nernst Equation (Non-Standard Conditions):
$$E = E^\circ - \frac{0.0592}{n} \log Q \quad (\text{at } 298\text{ K})$$`,
						PdfUrl: "sample.pdf",
					},
				},
			},
		},
		Quizzes: []QuizData{
			{
				Title:            "Chemistry Quiz 1: Quantum Numbers & Atomic Bonding",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "Which set of quantum numbers (n, l, ml, ms) is NOT allowed for an electron?",
						Explanation: "The angular momentum quantum number l can only range from 0 to (n-1). For n = 2, l cannot be 2.",
						Options: []Option{
							{Text: "n = 2, l = 2, ml = 0, ms = +1/2", IsCorrect: true},
							{Text: "n = 3, l = 1, ml = -1, ms = -1/2", IsCorrect: false},
							{Text: "n = 1, l = 0, ml = 0, ms = +1/2", IsCorrect: false},
							{Text: "n = 4, l = 3, ml = +2, ms = -1/2", IsCorrect: false},
						},
					},
					{
						Text:        "What is the molecular geometry of a water molecule (H2O) according to VSEPR theory?",
						Explanation: "H2O has 2 bonding pairs and 2 lone pairs on oxygen, giving a bent (angular) shape with bond angle ~104.5 deg.",
						Options: []Option{
							{Text: "Bent / Angular", IsCorrect: true},
							{Text: "Linear", IsCorrect: false},
							{Text: "Tetrahedral", IsCorrect: false},
							{Text: "Trigonal planar", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Chemistry Quiz 2: Kinetics & Rate Laws",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "If doubling the concentration of reactant A quadruples the initial reaction rate, what is the order with respect to A?",
						Explanation: "Rate proportional to [A]^m => 2^m = 4 => m = 2 (Second order).",
						Options: []Option{
							{Text: "Second order", IsCorrect: true},
							{Text: "First order", IsCorrect: false},
							{Text: "Zero order", IsCorrect: false},
							{Text: "Third order", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Chemistry Quiz 3: Equilibrium & Le Chatelier",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "For the Haber process N2(g) + 3H2(g) <=> 2NH3(g) (Delta H = -92 kJ/mol), increasing pressure shifts equilibrium:",
						Explanation: "Increasing pressure favors the side with fewer gas moles (4 moles on left vs 2 moles on right), shifting equilibrium to the right.",
						Options: []Option{
							{Text: "To the right (producing more NH3)", IsCorrect: true},
							{Text: "To the left (producing more N2 and H2)", IsCorrect: false},
							{Text: "No effect on position", IsCorrect: false},
							{Text: "Decreases the value of Kc", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Chemistry Quiz 4: Acid-Base pH & Buffers",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "What is the pH of a 0.001 M HCl solution?",
						Explanation: "HCl is a strong acid, so [H+] = 0.001 M = 10^-3 M. pH = -log(10^-3) = 3.0.",
						Options: []Option{
							{Text: "3.0", IsCorrect: true},
							{Text: "1.0", IsCorrect: false},
							{Text: "11.0", IsCorrect: false},
							{Text: "4.0", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Chemistry Quiz 5: Electrochemistry & Galvanic Cells",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "In a standard Daniel galvanic cell (Zn-Cu), what reaction occurs at the anode?",
						Explanation: "Oxidation always occurs at the anode: Zn(s) -> Zn2+(aq) + 2e-.",
						Options: []Option{
							{Text: "Oxidation of Zinc: Zn(s) -> Zn2+(aq) + 2e-", IsCorrect: true},
							{Text: "Reduction of Copper: Cu2+(aq) + 2e- -> Cu(s)", IsCorrect: false},
							{Text: "Reduction of Zinc: Zn2+(aq) + 2e- -> Zn(s)", IsCorrect: false},
							{Text: "Oxidation of Copper: Cu(s) -> Cu2+(aq) + 2e-", IsCorrect: false},
						},
					},
				},
			},
		},
		Exams: []ExamData{
			{
				Title:           "Chemistry Unit 1 & 2 Physical Chemistry Exam",
				Subject:         "Chemistry",
				Grade:           12,
				DurationMinutes: 45,
				PassMarks:       12,
				Instructions:    "Covers Quantum Mechanics, Periodic Trends, Chemical Kinetics, and Equilibrium. 45 minutes.",
				Questions: []QuestionData{
					{
						Text:        "Which element has the highest first ionization energy among the following?",
						Explanation: "Helium (He) has a full 1s2 shell close to the nucleus and has the highest first ionization energy on the periodic table.",
						Options: []Option{
							{Text: "Helium", IsCorrect: true},
							{Text: "Fluorine", IsCorrect: false},
							{Text: "Neon", IsCorrect: false},
							{Text: "Oxygen", IsCorrect: false},
						},
					},
					{
						Text:        "If the rate constant k has units L/(mol*s), what is the overall reaction order?",
						Explanation: "Units of k = (mol/L)^(1-n) s^-1 = L^(n-1) mol^(1-n) s^-1. For n = 2, units are L mol^-1 s^-1 (Second order).",
						Options: []Option{
							{Text: "Second order", IsCorrect: true},
							{Text: "First order", IsCorrect: false},
							{Text: "Zero order", IsCorrect: false},
							{Text: "Third order", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:           "Chemistry Midterm Examination (Acids, Bases & Electrochemistry)",
				Subject:         "Chemistry",
				Grade:           12,
				DurationMinutes: 60,
				PassMarks:       15,
				Instructions:    "Testing pH Calculations, Buffer Solutions, Galvanic Cells, and Electrolysis.",
				Questions: []QuestionData{
					{
						Text:        "Given E°(Cu2+/Cu) = +0.34 V and E°(Zn2+/Zn) = -0.76 V, calculate the standard cell potential for the Daniel cell.",
						Explanation: "E°cell = E°cathode - E°anode = 0.34 - (-0.76) = 0.34 + 0.76 = +1.10 V.",
						Options: []Option{
							{Text: "+1.10 V", IsCorrect: true},
							{Text: "+0.42 V", IsCorrect: false},
							{Text: "-1.10 V", IsCorrect: false},
							{Text: "+0.76 V", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:           "Ethiopian National EUEE Chemistry Final Mock Exam",
				Subject:         "Chemistry",
				Grade:           12,
				DurationMinutes: 90,
				PassMarks:       25,
				Instructions:    "Comprehensive National Entrance Examination simulation for Grade 12 Natural Science Chemistry.",
				Questions: []QuestionData{
					{
						Text:        "How many Coulombs of charge are required to deposit 1 mole of Aluminum from molten Al2O3? (Al3+ + 3e- -> Al)",
						Explanation: "3 moles of electrons are required per mole of Al. Q = n * F = 3 * 96,500 = 289,500 C.",
						Options: []Option{
							{Text: "289,500 C", IsCorrect: true},
							{Text: "96,500 C", IsCorrect: false},
							{Text: "193,000 C", IsCorrect: false},
							{Text: "32,167 C", IsCorrect: false},
						},
					},
				},
			},
		},
	}
}
