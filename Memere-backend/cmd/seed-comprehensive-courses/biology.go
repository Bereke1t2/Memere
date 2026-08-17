package main

func buildBiologyCourse() CourseData {
	return CourseData{
		Title:            "Grade 12 Biology: Ethiopian National University Entrance Prep",
		Slug:             "grade-12-biology-euee-prep",
		Subject:          "Biology",
		Grade:            12,
		Description:      "In-depth Grade 11 & 12 Biology curriculum designed for high scores on the Ethiopian University Entrance Exam (EUEE). Master Cell Biology & Biochemistry, Molecular Genetics & DNA Replication, Human Physiology, Evolution & Darwinian Selection, and Ecology with practice quizzes and mock entrance tests.",
		ShortDescription: "Complete Grade 12 Biology EUEE prep with cell biochemistry, genetics, physiology, ecology, quizzes, and exams.",
		ThumbnailUrl:     "https://images.unsplash.com/photo-1530210124550-912dc1381cb8?w=800&auto=format&fit=crop&q=80",
		Price:            0,
		IsFree:           true,
		Level:            "advanced",
		Sections: []SectionData{
			{
				Title:       "Unit 1: Biochemical Molecules & Cellular Energetics",
				Description: "Carbohydrates, lipids, proteins, enzymes, ATP synthesis, glycolysis, Krebs cycle, and oxidative phosphorylation.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 1.1: Macromolecules & Enzyme Kinetics",
						Type:            "video",
						DurationSeconds: 1200,
						IsFreePreview:   true,
						Content: `# Biological Macromolecules and Enzymes
Life processes depend on four major classes of organic biomolecules:

### 1. Carbohydrates:
- Monosaccharides (Glucose, Fructose, Galactose) linked by **glycosidic bonds**.
- Storage: Glycogen (animals) and Starch (plants).

### 2. Proteins:
- Polymers of amino acids linked by **peptide bonds**.
- Primary, secondary ($\alpha$-helix, $\beta$-sheet), tertiary, and quaternary structures.

### 3. Enzymes:
- Biological catalysts that lower the **activation energy** ($E_a$) of reactions without being consumed.
- Induced-fit model: active site conformational change upon substrate binding.`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 1.2: Cellular Respiration & ATP Yield",
						Type:            "note",
						DurationSeconds: 950,
						IsFreePreview:   true,
						Content: `# Cellular Respiration Pathways
Aerobic cellular respiration converts glucose into usable ATP energy:
$$\text{C}_6\text{H}_{12}\text{O}_6 + 6\text{O}_2 \to 6\text{CO}_2 + 6\text{H}_2\text{O} + 36\text{-}38\text{ ATP}$$

### Four Stages:
1. **Glycolysis** (Cytoplasm): Glucose $\to$ 2 Pyruvate + 2 Net ATP + 2 NADH (Anaerobic).
2. **Link Reaction** (Mitochondrial Matrix): Pyruvate $\to$ Acetyl-CoA + $\text{CO}_2$ + NADH.
3. **Krebs Cycle** (Matrix): 2 ATP + 6 NADH + 2 $\text{FADH}_2$.
4. **Oxidative Phosphorylation** (Inner Membrane / Cristae): Electron Transport Chain + ATP Synthase via Chemiosmosis.`,
						PdfUrl: "sample.pdf",
					},
				},
			},
			{
				Title:       "Unit 2: Molecular Genetics & Protein Synthesis",
				Description: "DNA structure, semi-conservative replication, transcription, translation, and genetic mutations.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 2.1: DNA Structure & Semi-Conservative Replication",
						Type:            "video",
						DurationSeconds: 1350,
						IsFreePreview:   false,
						Content: `# DNA Structure & Replication
- Double-helix discovered by Watson & Crick with antiparallel strands ($5' \to 3'$ and $3' \to 5'$).
- Complementary base pairing: Adenine with Thymine (2 H-bonds), Guanine with Cytosine (3 H-bonds).

### Key Enzymes:
- **Helicase**: Unwinds the double helix at the replication fork.
- **DNA Polymerase III**: Synthesizes new strand in the $5' \to 3'$ direction.
- **Ligase**: Seals Okazaki fragments on the lagging strand.`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 2.2: Transcription, Translation & the Genetic Code",
						Type:            "note",
						DurationSeconds: 900,
						IsFreePreview:   false,
						Content: `# The Central Dogma of Molecular Biology
$$\text{DNA} \xrightarrow{\text{Transcription}} \text{mRNA} \xrightarrow{\text{Translation}} \text{Polypeptide (Protein)}$$

### Key Features of Genetic Code:
- **Triplet Code**: 3 nucleotide bases form a codon specifying 1 amino acid.
- **Universal & Degenerate**: Multiple codons can code for the same amino acid.
- **Start Codon**: AUG (Methionine).
- **Stop Codons**: UAA, UAG, UGA.`,
						PdfUrl: "sample.pdf",
					},
				},
			},
			{
				Title:       "Unit 3: Human Physiology & Homeostasis",
				Description: "Nervous impulse transmission, endocrine regulation, kidney nephron function, and immune defenses.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 3.1: Action Potentials & Synaptic Transmission",
						Type:            "video",
						DurationSeconds: 1400,
						IsFreePreview:   false,
						Content: `# Neural Coordination & Action Potential
- **Resting Membrane Potential**: $-70\text{ mV}$, maintained by the $\text{Na}^+/\text{K}^+$ pump ($3\text{Na}^+$ out for $2\text{K}^+$ in).
- **Depolarization**: Voltage-gated $\text{Na}^+$ channels open $\to$ influx of $\text{Na}^+$ makes inside $+30\text{ mV}$.
- **Repolarization**: $\text{K}^+$ channels open $\to$ efflux of $\text{K}^+$.`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 3.2: Excretion & Nephron Countercurrent Mechanism",
						Type:            "note",
						DurationSeconds: 850,
						IsFreePreview:   false,
						Content: `# Renal Physiology & Nephron Function
1. **Ultrafiltration**: In Bowman's capsule under high hydrostatic pressure from afferent arteriole.
2. **Selective Reabsorption**: In the Proximal Convoluted Tubule (PCT) - 100% of glucose & amino acids, 80% water.
3. **Loop of Henle**: Creates hypertonic medulla for water conservation.
4. **ADH (Vasopressin)**: Increases water permeability of collecting ducts.`,
						PdfUrl: "sample.pdf",
					},
				},
			},
			{
				Title:       "Unit 4: Evolution and Ecology",
				Description: "Evidence for evolution, natural selection mechanisms, population genetics (Hardy-Weinberg), and ecosystem energy flow.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 4.1: Mechanisms of Evolution & Speciation",
						Type:            "video",
						DurationSeconds: 1100,
						IsFreePreview:   false,
						Content: `# Darwinian Natural Selection
- Overproduction of offspring leads to competition for limited resources.
- Inherited phenotypic variation leads to differential reproductive success.

### Speciation Types:
- **Allopatric**: Speciation due to geographic isolation.
- **Sympatric**: Speciation within the same geographical range (e.g. polyploidy in plants).`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 4.2: Ecosystem Dynamics & Biogeochemical Cycles",
						Type:            "note",
						DurationSeconds: 900,
						IsFreePreview:   false,
						Content: `# Ecology & Energy Flow
- **10% Energy Rule**: Only about 10% of energy is transferred from one trophic level to the next.
- **Nitrogen Cycle**:
  - Nitrogen fixation: *Rhizobium* in legume root nodules.
  - Nitrification: *Nitrosomonas* ($\text{NH}_4^+ \to \text{NO}_2^-$) and *Nitrobacter* ($\text{NO}_2^- \to \text{NO}_3^-$).
  - Denitrification: *Pseudomonas* converts nitrates back to $\text{N}_2$ gas.`,
						PdfUrl: "sample.pdf",
					},
				},
			},
		},
		Quizzes: []QuizData{
			{
				Title:            "Biology Quiz 1: Cell Biochemistry & Enzymes",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "How do biological enzymes increase the rate of chemical reactions?",
						Explanation: "Enzymes increase the reaction rate by lowering the activation energy barrier required for the transition state.",
						Options: []Option{
							{Text: "By lowering the activation energy", IsCorrect: true},
							{Text: "By increasing the temperature", IsCorrect: false},
							{Text: "By altering the reaction equilibrium", IsCorrect: false},
							{Text: "By increasing the free energy of products", IsCorrect: false},
						},
					},
					{
						Text:        "Which cellular respiration process occurs in the cytoplasm and does not require oxygen?",
						Explanation: "Glycolysis occurs in the cytoplasm and is an anaerobic process producing 2 net ATP per glucose.",
						Options: []Option{
							{Text: "Glycolysis", IsCorrect: true},
							{Text: "Krebs cycle", IsCorrect: false},
							{Text: "Electron transport chain", IsCorrect: false},
							{Text: "Oxidative phosphorylation", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Biology Quiz 2: Molecular Genetics & DNA Replication",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "Which enzyme is responsible for sealing Okazaki fragments on the lagging strand during DNA replication?",
						Explanation: "DNA ligase joins Okazaki fragments by catalyzing the formation of phosphodiester bonds.",
						Options: []Option{
							{Text: "DNA ligase", IsCorrect: true},
							{Text: "DNA helicase", IsCorrect: false},
							{Text: "DNA polymerase I", IsCorrect: false},
							{Text: "RNA primase", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Biology Quiz 3: Protein Synthesis & Translation",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "Which codon serves as the universal start codon for translation in eukaryotic cells?",
						Explanation: "AUG is the universal start codon, coding for the amino acid Methionine.",
						Options: []Option{
							{Text: "AUG", IsCorrect: true},
							{Text: "UAA", IsCorrect: false},
							{Text: "UAG", IsCorrect: false},
							{Text: "UGA", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Biology Quiz 4: Nervous & Endocrine Coordination",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "During an action potential, rapid depolarization of the axon membrane is caused by the sudden influx of:",
						Explanation: "Opening of voltage-gated Sodium (Na+) channels allows Na+ ions to flood into the axon down their electrochemical gradient.",
						Options: []Option{
							{Text: "Sodium (Na+) ions", IsCorrect: true},
							{Text: "Potassium (K+) ions", IsCorrect: false},
							{Text: "Calcium (Ca2+) ions", IsCorrect: false},
							{Text: "Chloride (Cl-) ions", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Biology Quiz 5: Ecology & Biogeochemical Cycles",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "Which genus of bacteria converts nitrites (NO2-) into nitrates (NO3-) during the nitrogen cycle in soil?",
						Explanation: "Nitrobacter oxidizes nitrites into nitrates, the chemical form most readily absorbed by plant roots.",
						Options: []Option{
							{Text: "Nitrobacter", IsCorrect: true},
							{Text: "Nitrosomonas", IsCorrect: false},
							{Text: "Rhizobium", IsCorrect: false},
							{Text: "Pseudomonas", IsCorrect: false},
						},
					},
				},
			},
		},
		Exams: []ExamData{
			{
				Title:           "Biology Unit 1 & 2 Biochemistry & Genetics Exam",
				Subject:         "Biology",
				Grade:           12,
				DurationMinutes: 45,
				PassMarks:       12,
				Instructions:    "Testing Biomolecules, Cellular Respiration, DNA Replication, and Protein Synthesis. 45 minutes.",
				Questions: []QuestionData{
					{
						Text:        "Where in the eukaryotic cell does the Krebs (Citric Acid) cycle take place?",
						Explanation: "The Krebs cycle takes place in the fluid matrix of the mitochondria.",
						Options: []Option{
							{Text: "Mitochondrial matrix", IsCorrect: true},
							{Text: "Mitochondrial cristae", IsCorrect: false},
							{Text: "Cytoplasm", IsCorrect: false},
							{Text: "Nucleolus", IsCorrect: false},
						},
					},
					{
						Text:        "If a DNA template strand has sequence 3'-TAC GGA CTT-5', what is the complementary mRNA transcript?",
						Explanation: "Complementary mRNA base pairs: 3'-TAC GGA CTT-5' => 5'-AUG CCU GAA-3'.",
						Options: []Option{
							{Text: "5'-AUG CCU GAA-3'", IsCorrect: true},
							{Text: "5'-AUG CCT GAA-3'", IsCorrect: false},
							{Text: "5'-UAC GGA CUU-3'", IsCorrect: false},
							{Text: "5'-ATG CCT GAA-3'", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:           "Biology Midterm Examination (Physiology & Homeostasis)",
				Subject:         "Biology",
				Grade:           12,
				DurationMinutes: 60,
				PassMarks:       15,
				Instructions:    "Testing Nervous System, Nephron Excretion, Hormone Regulation, and Homeostasis.",
				Questions: []QuestionData{
					{
						Text:        "Which hormone stimulates the insertion of aquaporin water channels into the kidney collecting ducts to concentrate urine?",
						Explanation: "Antidiuretic Hormone (ADH / Vasopressin), released from the posterior pituitary, increases water reabsorption.",
						Options: []Option{
							{Text: "Antidiuretic Hormone (ADH)", IsCorrect: true},
							{Text: "Aldosterone", IsCorrect: false},
							{Text: "Insulin", IsCorrect: false},
							{Text: "Thyroxine", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:           "Ethiopian National EUEE Biology Final Mock Exam",
				Subject:         "Biology",
				Grade:           12,
				DurationMinutes: 90,
				PassMarks:       25,
				Instructions:    "Official National University Entrance Examination Simulation for Grade 12 Natural Science Biology.",
				Questions: []QuestionData{
					{
						Text:        "According to the 10 percent rule in ecological energy transfer, if producers contain 50,000 kJ of energy, how much energy is available to secondary consumers?",
						Explanation: "Producers (50,000 kJ) -> Primary Consumers (5,000 kJ) -> Secondary Consumers (500 kJ).",
						Options: []Option{
							{Text: "500 kJ", IsCorrect: true},
							{Text: "5,000 kJ", IsCorrect: false},
							{Text: "50 kJ", IsCorrect: false},
							{Text: "5 kJ", IsCorrect: false},
						},
					},
				},
			},
		},
	}
}
