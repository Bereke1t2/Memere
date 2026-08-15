package main

import "fmt"

func buildBiologyExam() ExamData {
	questions := []QuestionData{
		{
			Text: "Which organelle is responsible for cellular respiration and ATP synthesis in eukaryotic cells?",
			Topic: "Cell Biology",
			Explanation: "Mitochondria are the powerhouses of the cell where the Krebs cycle and oxidative phosphorylation occur.",
			Options: []Option{
				{Text: "Mitochondrion", IsCorrect: true},
				{Text: "Ribosome", IsCorrect: false},
				{Text: "Golgi apparatus", IsCorrect: false},
				{Text: "Endoplasmic reticulum", IsCorrect: false},
			},
		},
		{
			Text: "What is the primary structural component of plant cell walls?",
			Topic: "Cell Biology",
			Explanation: "Cellulose, a complex polysaccharide of glucose units, provides rigidity and strength to plant cell walls.",
			Options: []Option{
				{Text: "Cellulose", IsCorrect: true},
				{Text: "Chitin", IsCorrect: false},
				{Text: "Peptidoglycan", IsCorrect: false},
				{Text: "Glycogen", IsCorrect: false},
			},
		},
		{
			Text: "Which phase of mitosis is characterized by chromosomes aligning along the equatorial plane of the cell?",
			Topic: "Cell Division",
			Explanation: "During metaphase, duplicated chromosomes align along the metaphase plate before sister chromatids separate.",
			Options: []Option{
				{Text: "Metaphase", IsCorrect: true},
				{Text: "Prophase", IsCorrect: false},
				{Text: "Anaphase", IsCorrect: false},
				{Text: "Telophase", IsCorrect: false},
			},
		},
		{
			Text: "In DNA structure, which nitrogenous base pairs with Adenine through two hydrogen bonds?",
			Topic: "Genetics & Molecular Biology",
			Explanation: "Adenine (A) pairs with Thymine (T) via 2 hydrogen bonds, while Guanine (G) pairs with Cytosine (C) via 3.",
			Options: []Option{
				{Text: "Thymine", IsCorrect: true},
				{Text: "Cytosine", IsCorrect: false},
				{Text: "Guanine", IsCorrect: false},
				{Text: "Uracil", IsCorrect: false},
			},
		},
		{
			Text: "What is the process of synthesizing mRNA from a DNA template called?",
			Topic: "Genetics & Molecular Biology",
			Explanation: "Transcription is the synthesis of RNA using a DNA template catalyzed by RNA polymerase.",
			Options: []Option{
				{Text: "Transcription", IsCorrect: true},
				{Text: "Translation", IsCorrect: false},
				{Text: "Replication", IsCorrect: false},
				{Text: "Transduction", IsCorrect: false},
			},
		},
		{
			Text: "According to Mendel's law of segregation, what phenotypic ratio is expected in the F2 generation of a monohybrid cross with complete dominance?",
			Topic: "Genetics & Heredity",
			Explanation: "A standard monohybrid cross (Bb x Bb) produces a 3:1 phenotypic ratio (3 dominant : 1 recessive).",
			Options: []Option{
				{Text: "3 : 1", IsCorrect: true},
				{Text: "9 : 3 : 3 : 1", IsCorrect: false},
				{Text: "1 : 2 : 1", IsCorrect: false},
				{Text: "1 : 1", IsCorrect: false},
			},
		},
		{
			Text: "Which enzyme is responsible for breaking down starch into maltose in human saliva?",
			Topic: "Human Physiology - Digestion",
			Explanation: "Salivary amylase (ptyalin) initiates starch digestion in the oral cavity.",
			Options: []Option{
				{Text: "Salivary amylase", IsCorrect: true},
				{Text: "Pepsin", IsCorrect: false},
				{Text: "Trypsin", IsCorrect: false},
				{Text: "Lipase", IsCorrect: false},
			},
		},
		{
			Text: "Which blood vessel carries oxygenated blood from the lungs back to the left atrium of the heart?",
			Topic: "Human Physiology - Circulation",
			Explanation: "The pulmonary veins are unique among veins because they carry oxygen-rich blood from lungs to heart.",
			Options: []Option{
				{Text: "Pulmonary vein", IsCorrect: true},
				{Text: "Pulmonary artery", IsCorrect: false},
				{Text: "Aorta", IsCorrect: false},
				{Text: "Superior vena cava", IsCorrect: false},
			},
		},
		{
			Text: "What is the functional unit of the human kidney responsible for filtration and urine formation?",
			Topic: "Human Physiology - Excretion",
			Explanation: "The nephron is the microscopic structural and functional unit of the kidney.",
			Options: []Option{
				{Text: "Nephron", IsCorrect: true},
				{Text: "Neuron", IsCorrect: false},
				{Text: "Alveolus", IsCorrect: false},
				{Text: "Glomerulus capsule", IsCorrect: false},
			},
		},
		{
			Text: "Which hormone produced by beta cells of the islets of Langerhans lowers blood glucose concentration?",
			Topic: "Endocrine System",
			Explanation: "Insulin promotes glucose uptake into cells and conversion of glucose to glycogen, lowering blood glucose levels.",
			Options: []Option{
				{Text: "Insulin", IsCorrect: true},
				{Text: "Glucagon", IsCorrect: false},
				{Text: "Adrenaline", IsCorrect: false},
				{Text: "Thyroxine", IsCorrect: false},
			},
		},
		{
			Text: "What is the primary site of gaseous exchange in human lungs?",
			Topic: "Respiratory System",
			Explanation: "Alveoli provide a huge, thin, moist surface area surrounded by capillaries for gas diffusion.",
			Options: []Option{
				{Text: "Alveoli", IsCorrect: true},
				{Text: "Bronchioles", IsCorrect: false},
				{Text: "Trachea", IsCorrect: false},
				{Text: "Pleural cavity", IsCorrect: false},
			},
		},
		{
			Text: "Which pigment in chloroplasts is primarily responsible for absorbing red and blue light during photosynthesis?",
			Topic: "Plant Physiology - Photosynthesis",
			Explanation: "Chlorophyll a and b absorb blue and red light while reflecting green light.",
			Options: []Option{
				{Text: "Chlorophyll", IsCorrect: true},
				{Text: "Carotenoid", IsCorrect: false},
				{Text: "Anthocyanin", IsCorrect: false},
				{Text: "Xanthophyll", IsCorrect: false},
			},
		},
		{
			Text: "In plants, water and dissolved mineral nutrients are transported upward from roots to leaves through:",
			Topic: "Plant Physiology - Transport",
			Explanation: "Xylem vessels transport water and dissolved inorganic ions driven by transpiration pull.",
			Options: []Option{
				{Text: "Xylem", IsCorrect: true},
				{Text: "Phloem", IsCorrect: false},
				{Text: "Cortex", IsCorrect: false},
				{Text: "Pith", IsCorrect: false},
			},
		},
		{
			Text: "Which plant hormone promotes cell elongation, apical dominance, and phototropism?",
			Topic: "Plant Hormones",
			Explanation: "Auxin (IAA) stimulates cell elongation and directs growth towards light.",
			Options: []Option{
				{Text: "Auxin", IsCorrect: true},
				{Text: "Ethylene", IsCorrect: false},
				{Text: "Abscisic acid", IsCorrect: false},
				{Text: "Cytokinin", IsCorrect: false},
			},
		},
		{
			Text: "What type of ecological relationship exists between mycorrhizal fungi and plant roots where both organisms benefit?",
			Topic: "Ecology",
			Explanation: "Mutualism is a symbiotic relationship in which both participating species derive a fitness benefit.",
			Options: []Option{
				{Text: "Mutualism", IsCorrect: true},
				{Text: "Commensalism", IsCorrect: false},
				{Text: "Parasitism", IsCorrect: false},
				{Text: "Amensalism", IsCorrect: false},
			},
		},
		{
			Text: "In an ecological food chain, approximately what percentage of energy is transferred from one trophic level to the next?",
			Topic: "Ecology - Energy Flow",
			Explanation: "According to Lindeman's 10% law, about 10% of the energy is transferred to the next trophic level.",
			Options: []Option{
				{Text: "10%", IsCorrect: true},
				{Text: "50%", IsCorrect: false},
				{Text: "1%", IsCorrect: false},
				{Text: "90%", IsCorrect: false},
			},
		},
		{
			Text: "Which evolutionary mechanism refers to random fluctuations in allele frequencies in small populations?",
			Topic: "Evolutionary Biology",
			Explanation: "Genetic drift is the change in allele frequency due to chance events, especially influential in small populations.",
			Options: []Option{
				{Text: "Genetic drift", IsCorrect: true},
				{Text: "Natural selection", IsCorrect: false},
				{Text: "Gene flow", IsCorrect: false},
				{Text: "Adaptive radiation", IsCorrect: false},
			},
		},
		{
			Text: "Homologous structures, such as the forelimbs of humans, bats, and whales, provide evidence of:",
			Topic: "Evolutionary Biology",
			Explanation: "Homologous structures share a common anatomical origin, demonstrating divergent evolution from a common ancestor.",
			Options: []Option{
				{Text: "Common ancestry and divergent evolution", IsCorrect: true},
				{Text: "Convergent evolution", IsCorrect: false},
				{Text: "Analogous adaptation", IsCorrect: false},
				{Text: "Spontaneous generation", IsCorrect: false},
			},
		},
		{
			Text: "Which technique is used in biotechnology to amplify specific sequences of DNA in vitro?",
			Topic: "Biotechnology",
			Explanation: "Polymerase Chain Reaction (PCR) enables exponential amplification of target DNA fragments.",
			Options: []Option{
				{Text: "Polymerase Chain Reaction (PCR)", IsCorrect: true},
				{Text: "Gel electrophoresis only", IsCorrect: false},
				{Text: "ELISA", IsCorrect: false},
				{Text: "Spectrophotometry", IsCorrect: false},
			},
		},
		{
			Text: "What type of pathogen causes human malaria?",
			Topic: "Microbiology & Diseases",
			Explanation: "Malaria is caused by the protozoan parasite of the genus Plasmodium (e.g., P. falciparum), transmitted by Anopheles mosquitoes.",
			Options: []Option{
				{Text: "Protozoan (Plasmodium)", IsCorrect: true},
				{Text: "Bacterium", IsCorrect: false},
				{Text: "Virus", IsCorrect: false},
				{Text: "Fungus", IsCorrect: false},
			},
		},
	}

	biologyTopics := []string{
		"Cellular Metabolism & Energetics", "Molecular Genetics & DNA Technology",
		"Mendelian & Human Genetics", "Organ Systems & Homeostasis",
		"Immunology & Disease", "Plant Growth & Responses",
		"Population Ecology & Biomes", "Evolution & Speciation",
	}

	for i := 21; i <= 60; i++ {
		topic := biologyTopics[i%len(biologyTopics)]
		var q QuestionData
		switch i % 8 {
		case 1:
			q = QuestionData{
				Text: fmt.Sprintf("During the light-dependent reactions of photosynthesis, photolysis of water produces oxygen, electrons, and:"),
				Topic: topic,
				Explanation: "Photolysis of H₂O splits water into 2H⁺ (protons), electrons (for PS II), and O₂ gas.",
				Options: []Option{
					{Text: "Protons (H⁺ ions)", IsCorrect: true},
					{Text: "Carbon dioxide (CO₂)", IsCorrect: false},
					{Text: "Glucose", IsCorrect: false},
					{Text: "ATP directly without synthase", IsCorrect: false},
				},
			}
		case 2:
			q = QuestionData{
				Text: fmt.Sprintf("How many net ATP molecules are yielded per glucose molecule during anaerobic glycolysis?"),
				Topic: topic,
				Explanation: "Glycolysis consumes 2 ATP and produces 4 ATP, yielding a net gain of 2 ATP per glucose.",
				Options: []Option{
					{Text: "2 ATP", IsCorrect: true},
					{Text: "36 ATP", IsCorrect: false},
					{Text: "4 ATP", IsCorrect: false},
					{Text: "38 ATP", IsCorrect: false},
				},
			}
		case 3:
			q = QuestionData{
				Text: fmt.Sprintf("Which type of RNA molecule carries amino acids to the ribosome during protein translation?"),
				Topic: topic,
				Explanation: "tRNA (transfer RNA) contains an anticodon and carries the corresponding specific amino acid to the ribosome.",
				Options: []Option{
					{Text: "Transfer RNA (tRNA)", IsCorrect: true},
					{Text: "Messenger RNA (mRNA)", IsCorrect: false},
					{Text: "Ribosomal RNA (rRNA)", IsCorrect: false},
					{Text: "Small nuclear RNA (snRNA)", IsCorrect: false},
				},
			}
		case 4:
			q = QuestionData{
				Text: fmt.Sprintf("Which blood type is considered the universal red blood cell donor because it lacks A and B surface antigens?"),
				Topic: topic,
				Explanation: "Type O-negative (O) red blood cells lack A and B agglutinogens and Rh antigen.",
				Options: []Option{
					{Text: "Type O", IsCorrect: true},
					{Text: "Type AB", IsCorrect: false},
					{Text: "Type A", IsCorrect: false},
					{Text: "Type B", IsCorrect: false},
				},
			}
		case 5:
			q = QuestionData{
				Text: fmt.Sprintf("In the human nervous system, the gap across which a nerve impulse passes from one neuron to another is the:"),
				Topic: topic,
				Explanation: "A synapse is the microscopic junction across which neurotransmitters diffuse from presynaptic to postsynaptic neurons.",
				Options: []Option{
					{Text: "Synapse", IsCorrect: true},
					{Text: "Axon hillock", IsCorrect: false},
					{Text: "Myelin sheath", IsCorrect: false},
					{Text: "Node of Ranvier", IsCorrect: false},
				},
			}
		case 6:
			q = QuestionData{
				Text: fmt.Sprintf("Which class of antibodies (immunoglobulins) is the most abundant in human blood serum and crosses the placenta?"),
				Topic: topic,
				Explanation: "IgG accounts for about 75-80% of all antibodies in the body and confers passive immunity to the fetus.",
				Options: []Option{
					{Text: "IgG", IsCorrect: true},
					{Text: "IgM", IsCorrect: false},
					{Text: "IgA", IsCorrect: false},
					{Text: "IgE", IsCorrect: false},
				},
			}
		case 7:
			q = QuestionData{
				Text: fmt.Sprintf("In double fertilization in angiosperms, one sperm fertilizes the egg while the second sperm fuses with two polar nuclei to form:"),
				Topic: topic,
				Explanation: "The triploid (3n) endosperm provides nutritive tissue for the developing embryo.",
				Options: []Option{
					{Text: "Triploid endosperm (3n)", IsCorrect: true},
					{Text: "Diploid zygote (2n)", IsCorrect: false},
					{Text: "Seed coat (testa)", IsCorrect: false},
					{Text: "Cotyledon only", IsCorrect: false},
				},
			}
		default:
			q = QuestionData{
				Text: fmt.Sprintf("In an ecosystem, organisms that break down dead organic matter and recycle vital nutrients are classified as:"),
				Topic: topic,
				Explanation: "Decomposers (saprotrophs such as fungi and bacteria) recycle carbon, nitrogen, and minerals back into the soil and atmosphere.",
				Options: []Option{
					{Text: "Decomposers (Saprotrophs)", IsCorrect: true},
					{Text: "Primary producers", IsCorrect: false},
					{Text: "Tertiary consumers", IsCorrect: false},
					{Text: "Herbivores", IsCorrect: false},
				},
			}
		}
		questions = append(questions, q)
	}

	return ExamData{
		Title:           "Grade 12 National Biology Mock Examination",
		Subject:         "Biology",
		Grade:           12,
		DurationMinutes: 120,
		PassMarks:       30,
		Instructions:    "This examination contains 60 multiple-choice questions covering Cell Biology, Genetics, Human & Plant Physiology, Ecology, and Evolution. Select the best answer for each question.",
		Questions:       questions,
	}
}
