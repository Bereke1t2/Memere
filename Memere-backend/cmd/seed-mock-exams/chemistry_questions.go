package main

import "fmt"

func buildChemistryExam() ExamData {
	questions := []QuestionData{
		{
			Text: "What is the electronic configuration of Calcium (atomic number Z = 20)?",
			Topic: "Atomic Structure",
			Explanation: "Calcium (Z=20) has 20 electrons: 1s² 2s² 2p⁶ 3s² 3p⁶ 4s² or [Ar] 4s².",
			Options: []Option{
				{Text: "1s² 2s² 2p⁶ 3s² 3p⁶ 4s²", IsCorrect: true},
				{Text: "1s² 2s² 2p⁶ 3s² 3p⁶ 3d²", IsCorrect: false},
				{Text: "1s² 2s² 2p⁶ 3s² 3d⁸", IsCorrect: false},
				{Text: "1s² 2s² 2p⁶ 3s² 3p⁴ 4s⁴", IsCorrect: false},
			},
		},
		{
			Text: "Which type of chemical bond is formed by the sharing of electron pairs between nonmetal atoms?",
			Topic: "Chemical Bonding",
			Explanation: "Covalent bonding involves the mutual sharing of valence electrons between atoms.",
			Options: []Option{
				{Text: "Covalent bond", IsCorrect: true},
				{Text: "Ionic bond", IsCorrect: false},
				{Text: "Metallic bond", IsCorrect: false},
				{Text: "Hydrogen bond", IsCorrect: false},
			},
		},
		{
			Text: "What is the molecular geometry of methane (CH₄) according to VSEPR theory?",
			Topic: "Chemical Bonding",
			Explanation: "CH₄ has 4 bonding pairs and 0 lone pairs around carbon, adopting a tetrahedral geometry with bond angles of 109.5°.",
			Options: []Option{
				{Text: "Tetrahedral", IsCorrect: true},
				{Text: "Trigonal planar", IsCorrect: false},
				{Text: "Linear", IsCorrect: false},
				{Text: "Octahedral", IsCorrect: false},
			},
		},
		{
			Text: "What is the oxidation state of sulfur in sulfuric acid (H₂SO₄)?",
			Topic: "Redox & Electrochemistry",
			Explanation: "2(+1) + S + 4(-2) = 0 => +2 + S - 8 = 0 => S = +6.",
			Options: []Option{
				{Text: "+6", IsCorrect: true},
				{Text: "+4", IsCorrect: false},
				{Text: "+2", IsCorrect: false},
				{Text: "-2", IsCorrect: false},
			},
		},
		{
			Text: "According to Le Chatelier's principle, what happens to the exothermic equilibrium N₂(g) + 3H₂(g) ⇌ 2NH₃(g) + heat when temperature is increased?",
			Topic: "Chemical Equilibrium",
			Explanation: "Increasing temperature shifts the equilibrium in the endothermic direction (reverse direction, toward reactants).",
			Options: []Option{
				{Text: "Shifts toward the reactants (left)", IsCorrect: true},
				{Text: "Shifts toward the products (right)", IsCorrect: false},
				{Text: "No change in equilibrium position", IsCorrect: false},
				{Text: "The equilibrium constant increases", IsCorrect: false},
			},
		},
		{
			Text: "What is the pH of a 0.001 M HCl solution at 25°C?",
			Topic: "Acid-Base Chemistry",
			Explanation: "HCl is a strong acid, [H⁺] = 10⁻³ M. pH = -log[H⁺] = -log(10⁻³) = 3.",
			Options: []Option{
				{Text: "3", IsCorrect: true},
				{Text: "1", IsCorrect: false},
				{Text: "11", IsCorrect: false},
				{Text: "7", IsCorrect: false},
			},
		},
		{
			Text: "Which of the following acts as a Lewis base?",
			Topic: "Acid-Base Chemistry",
			Explanation: "A Lewis base is an electron pair donor. Ammonia (NH₃) has a lone pair on nitrogen.",
			Options: []Option{
				{Text: "NH₃ (Ammonia)", IsCorrect: true},
				{Text: "BF₃ (Boron trifluoride)", IsCorrect: false},
				{Text: "AlCl₃ (Aluminum chloride)", IsCorrect: false},
				{Text: "H⁺ (Proton)", IsCorrect: false},
			},
		},
		{
			Text: "What is the IUPAC name for the organic compound CH₃-CH(CH₃)-CH₂-CH₃?",
			Topic: "Organic Chemistry",
			Explanation: "The longest continuous carbon chain has 4 carbons (butane) with a methyl group at position 2: 2-methylbutane.",
			Options: []Option{
				{Text: "2-methylbutane", IsCorrect: true},
				{Text: "3-methylbutane", IsCorrect: false},
				{Text: "Isopentane", IsCorrect: false},
				{Text: "Dimethylpropane", IsCorrect: false},
			},
		},
		{
			Text: "Which functional group is characteristic of carboxylic acids?",
			Topic: "Organic Chemistry",
			Explanation: "Carboxylic acids contain the carboxyl group (-COOH).",
			Options: []Option{
				{Text: "-COOH", IsCorrect: true},
				{Text: "-CHO", IsCorrect: false},
				{Text: "-OH", IsCorrect: false},
				{Text: "-CO-", IsCorrect: false},
			},
		},
		{
			Text: "What is the product of the complete combustion of propane (C₃H₈) in excess oxygen?",
			Topic: "Organic Chemistry",
			Explanation: "Complete combustion of hydrocarbons in excess O₂ yields carbon dioxide (CO₂) and water (H₂O).",
			Options: []Option{
				{Text: "CO₂ and H₂O", IsCorrect: true},
				{Text: "CO and H₂O", IsCorrect: false},
				{Text: "C and H₂O", IsCorrect: false},
				{Text: "CO₂ and H₂", IsCorrect: false},
			},
		},
		{
			Text: "Which law states that the rate of effusion of a gas is inversely proportional to the square root of its molar mass?",
			Topic: "States of Matter",
			Explanation: "Graham's Law of Effusion: Rate ∝ 1 / √(Molar Mass).",
			Options: []Option{
				{Text: "Graham's Law", IsCorrect: true},
				{Text: "Dalton's Law", IsCorrect: false},
				{Text: "Charles's Law", IsCorrect: false},
				{Text: "Boyle's Law", IsCorrect: false},
			},
		},
		{
			Text: "What is the standard enthalpy of formation (ΔH°f) for pure elements in their standard states?",
			Topic: "Thermodynamics",
			Explanation: "By convention, the standard enthalpy of formation of an element in its standard reference state is zero (0 kJ/mol).",
			Options: []Option{
				{Text: "0 kJ/mol", IsCorrect: true},
				{Text: "100 kJ/mol", IsCorrect: false},
				{Text: "-285.8 kJ/mol", IsCorrect: false},
				{Text: "1 kJ/mol", IsCorrect: false},
			},
		},
		{
			Text: "In an electrolytic cell, oxidation always takes place at which electrode?",
			Topic: "Electrochemistry",
			Explanation: "Oxidation always occurs at the anode (An Ox) in both galvanic and electrolytic cells.",
			Options: []Option{
				{Text: "Anode", IsCorrect: true},
				{Text: "Cathode", IsCorrect: false},
				{Text: "Salt bridge", IsCorrect: false},
				{Text: "Electrolyte solution", IsCorrect: false},
			},
		},
		{
			Text: "What is the unit of the rate constant (k) for a first-order chemical reaction?",
			Topic: "Chemical Kinetics",
			Explanation: "Rate = k[A] => k = Rate / [A] = (M/s) / M = s⁻¹.",
			Options: []Option{
				{Text: "s⁻¹", IsCorrect: true},
				{Text: "M⁻¹ · s⁻¹", IsCorrect: false},
				{Text: "M · s⁻¹", IsCorrect: false},
				{Text: "M⁻² · s⁻¹", IsCorrect: false},
			},
		},
		{
			Text: "Which catalyst is traditionally used in the industrial Haber-Bosch process for ammonia synthesis?",
			Topic: "Industrial Chemistry",
			Explanation: "Finely divided iron (Fe) promoted with potassium and aluminum oxides is the standard Haber catalyst.",
			Options: []Option{
				{Text: "Iron (Fe)", IsCorrect: true},
				{Text: "Vanadium pentoxide (V₂O₅)", IsCorrect: false},
				{Text: "Platinum (Pt)", IsCorrect: false},
				{Text: "Nickel (Ni)", IsCorrect: false},
			},
		},
		{
			Text: "Which catalyst is used in the Contact Process for the manufacture of sulfuric acid?",
			Topic: "Industrial Chemistry",
			Explanation: "Vanadium(V) oxide (V₂O₅) is used to catalyze the oxidation of SO₂ to SO₃.",
			Options: []Option{
				{Text: "Vanadium(V) oxide (V₂O₅)", IsCorrect: true},
				{Text: "Iron oxide (Fe₂O₃)", IsCorrect: false},
				{Text: "Manganese dioxide (MnO₂)", IsCorrect: false},
				{Text: "Copper (Cu)", IsCorrect: false},
			},
		},
		{
			Text: "Which polymer is formed through the addition polymerization of ethene (C₂H₄)?",
			Topic: "Polymers",
			Explanation: "Polymerization of ethene (ethylene) produces polyethene (polyethylene).",
			Options: []Option{
				{Text: "Polyethylene", IsCorrect: true},
				{Text: "Polyvinyl chloride (PVC)", IsCorrect: false},
				{Text: "Nylon 6,6", IsCorrect: false},
				{Text: "Polystyrene", IsCorrect: false},
			},
		},
		{
			Text: "What is the primary greenhouse gas produced by human industrial activities and fossil fuel combustion?",
			Topic: "Environmental Chemistry",
			Explanation: "Carbon dioxide (CO₂) is the primary greenhouse gas emitted through anthropogenic activities.",
			Options: []Option{
				{Text: "Carbon dioxide (CO₂)", IsCorrect: true},
				{Text: "Sulfur dioxide (SO₂)", IsCorrect: false},
				{Text: "Chlorofluorocarbons (CFCs)", IsCorrect: false},
				{Text: "Argon (Ar)", IsCorrect: false},
			},
		},
		{
			Text: "Which element has the highest electronegativity on the Pauling scale?",
			Topic: "Periodic Table & Periodic Trends",
			Explanation: "Fluorine (F) is the most electronegative element with a value of 3.98 (approximately 4.0).",
			Options: []Option{
				{Text: "Fluorine (F)", IsCorrect: true},
				{Text: "Oxygen (O)", IsCorrect: false},
				{Text: "Chlorine (Cl)", IsCorrect: false},
				{Text: "Francium (Fr)", IsCorrect: false},
			},
		},
		{
			Text: "What is the molar mass of glucose (C₆H₁₂O₆)? (C=12, H=1, O=16 g/mol)",
			Topic: "Stoichiometry",
			Explanation: "6(12) + 12(1) + 6(16) = 72 + 12 + 96 = 180 g/mol.",
			Options: []Option{
				{Text: "180 g/mol", IsCorrect: true},
				{Text: "342 g/mol", IsCorrect: false},
				{Text: "90 g/mol", IsCorrect: false},
				{Text: "160 g/mol", IsCorrect: false},
			},
		},
	}

	chemistryTopics := []string{
		"Atomic Structure & Periodic Trends", "Chemical Bonding & Molecular Shapes",
		"Thermodynamics & Thermochemistry", "Kinetics & Reaction Mechanisms",
		"Equilibrium & Acid-Base Reactions", "Electrochemistry & Batteries",
		"Organic Reactions & Mechanisms", "Biomolecules & Polymers",
	}

	for i := 21; i <= 60; i++ {
		topic := chemistryTopics[i%len(chemistryTopics)]
		var q QuestionData
		switch i % 8 {
		case 1:
			q = QuestionData{
				Text: fmt.Sprintf("How many moles of gas are present in a %d L container at STP (0°C, 1 atm)? (Molar volume = 22.4 L/mol)", (i%5+1)*224/10),
				Topic: topic,
				Explanation: fmt.Sprintf("n = V / 22.4 = %d / 22.4 = %d moles.", (i%5+1)*224/10, i%5+1),
				Options: []Option{
					{Text: fmt.Sprintf("%d moles", i%5+1), IsCorrect: true},
					{Text: fmt.Sprintf("%d moles", (i%5+1)*2), IsCorrect: false},
					{Text: fmt.Sprintf("%.1f moles", float64(i%5+1)*0.5), IsCorrect: false},
					{Text: "22.4 moles", IsCorrect: false},
				},
			}
		case 2:
			q = QuestionData{
				Text: fmt.Sprintf("What is the conjugate base of the acid H%dPO4?", i%2+2),
				Topic: topic,
				Explanation: "Removing one proton (H⁺) yields the conjugate base.",
				Options: []Option{
					{Text: "H₂PO₄⁻", IsCorrect: true},
					{Text: "PO₄³⁻", IsCorrect: false},
					{Text: "H₃PO₄", IsCorrect: false},
					{Text: "HPO₄²⁻", IsCorrect: false},
				},
			}
		case 3:
			q = QuestionData{
				Text: fmt.Sprintf("If the activation energy of a reaction is lowered by %d kJ/mol using a catalyst, the reaction rate:", (i%5+2)*10),
				Topic: topic,
				Explanation: "A catalyst lowers activation energy, exponentially increasing the fraction of effective collisions and thus increasing rate.",
				Options: []Option{
					{Text: "Increases significantly", IsCorrect: true},
					{Text: "Decreases significantly", IsCorrect: false},
					{Text: "Remains unchanged", IsCorrect: false},
					{Text: "Shifts equilibrium toward reactants", IsCorrect: false},
				},
			}
		case 4:
			q = QuestionData{
				Text: fmt.Sprintf("What is the hybridization of carbon in %s?", []string{"ethyne (HC≡CH)", "ethene (H₂C=CH₂)", "methane (CH₄)"}[i%3]),
				Topic: topic,
				Explanation: "Triple bonded carbon is sp, double bonded is sp², single bonded is sp³.",
				Options: []Option{
					{Text: []string{"sp", "sp²", "sp³"}[i%3], IsCorrect: true},
					{Text: "sp³d", IsCorrect: false},
					{Text: "dsp²", IsCorrect: false},
					{Text: []string{"sp³", "sp", "sp²"}[i%3], IsCorrect: false},
				},
			}
		case 5:
			q = QuestionData{
				Text: fmt.Sprintf("What mass of NaCl (molar mass = 58.5 g/mol) is required to prepare 500 mL of a %d.0 M solution?", i%3+1),
				Topic: topic,
				Explanation: fmt.Sprintf("Moles = M * V = %d.0 * 0.5 = %.1f mol. Mass = %.1f * 58.5 = %.1f g.", i%3+1, float64(i%3+1)*0.5, float64(i%3+1)*0.5, float64(i%3+1)*0.5*58.5),
				Options: []Option{
					{Text: fmt.Sprintf("%.1f g", float64(i%3+1)*0.5*58.5), IsCorrect: true},
					{Text: fmt.Sprintf("%.1f g", float64(i%3+1)*58.5), IsCorrect: false},
					{Text: "58.5 g", IsCorrect: false},
					{Text: "117.0 g", IsCorrect: false},
				},
			}
		case 6:
			q = QuestionData{
				Text: fmt.Sprintf("Which of the following organic reactions converts an alcohol to an alkene? (Example with %d carbons)", i%4+2),
				Topic: topic,
				Explanation: "Acid-catalyzed dehydration of alcohols removes a molecule of water to form an alkene (elimination reaction).",
				Options: []Option{
					{Text: "Dehydration (Elimination)", IsCorrect: true},
					{Text: "Hydrogenation (Addition)", IsCorrect: false},
					{Text: "Hydrolysis", IsCorrect: false},
					{Text: "Saponification", IsCorrect: false},
				},
			}
		case 7:
			q = QuestionData{
				Text: fmt.Sprintf("What is the standard cell potential E°cell for a cell with E°cathode = +0.%d V and E°anode = -0.%d V?", (i%4+5), (i%3+2)),
				Topic: topic,
				Explanation: fmt.Sprintf("E°cell = E°cathode - E°anode = 0.%d - (-0.%d) = %.2f V.", i%4+5, i%3+2, float64(i%4+5)/10.0+float64(i%3+2)/10.0),
				Options: []Option{
					{Text: fmt.Sprintf("%.2f V", float64(i%4+5)/10.0+float64(i%3+2)/10.0), IsCorrect: true},
					{Text: fmt.Sprintf("%.2f V", float64(i%4+5)/10.0-float64(i%3+2)/10.0), IsCorrect: false},
					{Text: "0.00 V", IsCorrect: false},
					{Text: "1.10 V", IsCorrect: false},
				},
			}
		default:
			q = QuestionData{
				Text: fmt.Sprintf("In the periodic table, atomic radius generally decreases across a period from left to right due to:"),
				Topic: topic,
				Explanation: "Increasing effective nuclear charge (Z_eff) pulls the valence electrons closer to the nucleus.",
				Options: []Option{
					{Text: "Increasing effective nuclear charge (Z_eff)", IsCorrect: true},
					{Text: "Increasing shielding effect", IsCorrect: false},
					{Text: "Decreasing principal quantum number", IsCorrect: false},
					{Text: "Decreasing number of protons", IsCorrect: false},
				},
			}
		}
		questions = append(questions, q)
	}

	return ExamData{
		Title:           "Grade 12 National Chemistry Mock Examination",
		Subject:         "Chemistry",
		Grade:           12,
		DurationMinutes: 120,
		PassMarks:       30,
		Instructions:    "This examination contains 60 multiple-choice questions covering General, Physical, Inorganic, and Organic Chemistry. Choose the correct option for each question.",
		Questions:       questions,
	}
}
