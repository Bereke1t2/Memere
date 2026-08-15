package main

import "fmt"

func buildPhysicsExam() ExamData {
	questions := []QuestionData{
		{
			Text: "A car accelerates uniformly from rest to a speed of 20 m/s in 5 seconds. What is its acceleration?",
			Topic: "Kinematics",
			Explanation: "a = (v - u) / t = (20 - 0) / 5 = 4 m/s².",
			Options: []Option{
				{Text: "4 m/s²", IsCorrect: true},
				{Text: "5 m/s²", IsCorrect: false},
				{Text: "2 m/s²", IsCorrect: false},
				{Text: "10 m/s²", IsCorrect: false},
			},
		},
		{
			Text: "A ball is thrown vertically upward with an initial velocity of 29.4 m/s. How long does it take to reach maximum height? (g = 9.8 m/s²)",
			Topic: "Kinematics",
			Explanation: "At maximum height v = 0. t = u / g = 29.4 / 9.8 = 3.0 s.",
			Options: []Option{
				{Text: "3.0 s", IsCorrect: true},
				{Text: "2.0 s", IsCorrect: false},
				{Text: "4.5 s", IsCorrect: false},
				{Text: "6.0 s", IsCorrect: false},
			},
		},
		{
			Text: "What is the magnitude of the net force required to give a 5 kg mass an acceleration of 6 m/s²?",
			Topic: "Dynamics & Newton's Laws",
			Explanation: "F = m * a = 5 kg * 6 m/s² = 30 N.",
			Options: []Option{
				{Text: "30 N", IsCorrect: true},
				{Text: "25 N", IsCorrect: false},
				{Text: "1.2 N", IsCorrect: false},
				{Text: "11 N", IsCorrect: false},
			},
		},
		{
			Text: "A 2 kg object moving at 4 m/s possesses how much kinetic energy?",
			Topic: "Work, Energy & Power",
			Explanation: "KE = 0.5 * m * v² = 0.5 * 2 * (4)² = 16 J.",
			Options: []Option{
				{Text: "16 J", IsCorrect: true},
				{Text: "8 J", IsCorrect: false},
				{Text: "32 J", IsCorrect: false},
				{Text: "4 J", IsCorrect: false},
			},
		},
		{
			Text: "How much work is done when a force of 50 N moves a box 4 meters in the direction of the force?",
			Topic: "Work, Energy & Power",
			Explanation: "W = F * d = 50 N * 4 m = 200 J.",
			Options: []Option{
				{Text: "200 J", IsCorrect: true},
				{Text: "12.5 J", IsCorrect: false},
				{Text: "54 J", IsCorrect: false},
				{Text: "250 J", IsCorrect: false},
			},
		},
		{
			Text: "What is the gravitational potential energy of a 10 kg block lifted 5 meters above ground? (g = 9.8 m/s²)",
			Topic: "Gravitation",
			Explanation: "PE = m * g * h = 10 * 9.8 * 5 = 490 J.",
			Options: []Option{
				{Text: "490 J", IsCorrect: true},
				{Text: "50 J", IsCorrect: false},
				{Text: "98 J", IsCorrect: false},
				{Text: "245 J", IsCorrect: false},
			},
		},
		{
			Text: "According to Newton's Law of Universal Gravitation, if the distance between two masses is doubled, the force between them becomes:",
			Topic: "Gravitation",
			Explanation: "F is proportional to 1/r². If r is doubled, F becomes (1/2)² = 1/4 of its original value.",
			Options: []Option{
				{Text: "One-fourth (1/4)", IsCorrect: true},
				{Text: "One-half (1/2)", IsCorrect: false},
				{Text: "Double (2x)", IsCorrect: false},
				{Text: "Four times (4x)", IsCorrect: false},
			},
		},
		{
			Text: "A hydraulic press works on the basis of which principle?",
			Topic: "Fluid Mechanics",
			Explanation: "Pascal's principle states that pressure applied to an enclosed fluid is transmitted undiminished in all directions.",
			Options: []Option{
				{Text: "Pascal's Principle", IsCorrect: true},
				{Text: "Archimedes' Principle", IsCorrect: false},
				{Text: "Bernoulli's Equation", IsCorrect: false},
				{Text: "Torricelli's Law", IsCorrect: false},
			},
		},
		{
			Text: "What is the buoyant force on an object displacing 0.002 m³ of water? (Density of water = 1000 kg/m³, g = 9.8 m/s²)",
			Topic: "Fluid Mechanics",
			Explanation: "F_b = ρ * V * g = 1000 * 0.002 * 9.8 = 19.6 N.",
			Options: []Option{
				{Text: "19.6 N", IsCorrect: true},
				{Text: "2.0 N", IsCorrect: false},
				{Text: "9.8 N", IsCorrect: false},
				{Text: "196 N", IsCorrect: false},
			},
		},
		{
			Text: "Which thermodynamic process occurs at constant pressure?",
			Topic: "Thermodynamics",
			Explanation: "An isobaric process is a thermodynamic process in which the pressure remains constant (ΔP = 0).",
			Options: []Option{
				{Text: "Isobaric", IsCorrect: true},
				{Text: "Isochoric", IsCorrect: false},
				{Text: "Isothermal", IsCorrect: false},
				{Text: "Adiabatic", IsCorrect: false},
			},
		},
		{
			Text: "What is the speed of sound in air at 20°C approximately?",
			Topic: "Waves & Oscillations",
			Explanation: "Speed of sound in dry air at 20°C is approximately 343 m/s.",
			Options: []Option{
				{Text: "343 m/s", IsCorrect: true},
				{Text: "300,000 km/s", IsCorrect: false},
				{Text: "1500 m/s", IsCorrect: false},
				{Text: "120 m/s", IsCorrect: false},
			},
		},
		{
			Text: "What is the frequency of a wave with speed 300 m/s and wavelength 1.5 meters?",
			Topic: "Waves & Oscillations",
			Explanation: "f = v / λ = 300 / 1.5 = 200 Hz.",
			Options: []Option{
				{Text: "200 Hz", IsCorrect: true},
				{Text: "450 Hz", IsCorrect: false},
				{Text: "150 Hz", IsCorrect: false},
				{Text: "0.005 Hz", IsCorrect: false},
			},
		},
		{
			Text: "What is the phenomenon responsible for the bending of light when it passes from air into water?",
			Topic: "Optics",
			Explanation: "Refraction is the change in direction and speed of a wave passing from one medium to another.",
			Options: []Option{
				{Text: "Refraction", IsCorrect: true},
				{Text: "Diffraction", IsCorrect: false},
				{Text: "Dispersion", IsCorrect: false},
				{Text: "Interference", IsCorrect: false},
			},
		},
		{
			Text: "Two point charges of +2 μC and +4 μC are separated by 0.2 m. If Coulomb constant k = 9 x 10⁹ N·m²/C², what is the electric force?",
			Topic: "Electrostatics",
			Explanation: "F = k*q1*q2 / r² = 9e9 * 2e-6 * 4e-6 / (0.04) = 72e-3 / 0.04 = 1.8 N.",
			Options: []Option{
				{Text: "1.8 N", IsCorrect: true},
				{Text: "3.6 N", IsCorrect: false},
				{Text: "0.9 N", IsCorrect: false},
				{Text: "18 N", IsCorrect: false},
			},
		},
		{
			Text: "According to Ohm's law, if a voltage of 12 V is applied across a 4 Ω resistor, what is the current?",
			Topic: "Current Electricity",
			Explanation: "I = V / R = 12 / 4 = 3 A.",
			Options: []Option{
				{Text: "3 A", IsCorrect: true},
				{Text: "48 A", IsCorrect: false},
				{Text: "0.33 A", IsCorrect: false},
				{Text: "8 A", IsCorrect: false},
			},
		},
		{
			Text: "Three 6 Ω resistors connected in parallel have an equivalent resistance of:",
			Topic: "Current Electricity",
			Explanation: "1/R_eq = 1/6 + 1/6 + 1/6 = 3/6 = 1/2 => R_eq = 2 Ω.",
			Options: []Option{
				{Text: "2 Ω", IsCorrect: true},
				{Text: "18 Ω", IsCorrect: false},
				{Text: "6 Ω", IsCorrect: false},
				{Text: "3 Ω", IsCorrect: false},
			},
		},
		{
			Text: "Faraday's law of electromagnetic induction states that the induced electromotive force (EMF) is proportional to:",
			Topic: "Electromagnetism",
			Explanation: "EMF = -dΦ/dt, proportional to the time rate of change of magnetic flux.",
			Options: []Option{
				{Text: "The rate of change of magnetic flux", IsCorrect: true},
				{Text: "The magnitude of the magnetic field only", IsCorrect: false},
				{Text: "The resistance of the coil only", IsCorrect: false},
				{Text: "The static electric charge", IsCorrect: false},
			},
		},
		{
			Text: "What is the energy of a photon of light with frequency 5 x 10¹⁴ Hz? (Planck's constant h = 6.63 x 10⁻³⁴ J·s)",
			Topic: "Modern & Quantum Physics",
			Explanation: "E = h * f = 6.63e-34 * 5e14 = 3.315 x 10⁻¹⁹ J.",
			Options: []Option{
				{Text: "3.315 × 10⁻¹⁹ J", IsCorrect: true},
				{Text: "1.326 × 10⁻⁴⁸ J", IsCorrect: false},
				{Text: "6.63 × 10⁻²⁰ J", IsCorrect: false},
				{Text: "7.54 × 10⁻¹⁹ J", IsCorrect: false},
			},
		},
		{
			Text: "In the photoelectric effect, the maximum kinetic energy of emitted electrons depends upon:",
			Topic: "Modern & Quantum Physics",
			Explanation: "KE_max = hf - Φ, which depends directly on the frequency of the incident light, not intensity.",
			Options: []Option{
				{Text: "Frequency of the incident radiation", IsCorrect: true},
				{Text: "Intensity of the incident radiation only", IsCorrect: false},
				{Text: "Duration of exposure only", IsCorrect: false},
				{Text: "Angle of incidence", IsCorrect: false},
			},
		},
		{
			Text: "What is the half-life of a radioactive isotope if 100 g decays to 12.5 g in 24 days?",
			Topic: "Nuclear Physics",
			Explanation: "100 -> 50 -> 25 -> 12.5 is 3 half-lives. 3 * T_half = 24 => T_half = 8 days.",
			Options: []Option{
				{Text: "8 days", IsCorrect: true},
				{Text: "6 days", IsCorrect: false},
				{Text: "12 days", IsCorrect: false},
				{Text: "4 days", IsCorrect: false},
			},
		},
	}

	physicsTopics := []string{
		"Mechanics & Rotational Motion", "Thermodynamics & Heat", "Oscillations & Waves",
		"Optics & Light", "Electrostatics & Capacitance", "Magnetism & Induction",
		"Alternating Current & Circuits", "Modern Physics & Relativity",
	}

	for i := 21; i <= 60; i++ {
		topic := physicsTopics[i%len(physicsTopics)]
		var q QuestionData
		switch i % 8 {
		case 1:
			q = QuestionData{
				Text: fmt.Sprintf("A wheel with radius 0.%d m rotates with an angular speed of 10 rad/s. What is the linear velocity of a point on the rim?", i%5+1),
				Topic: topic,
				Explanation: fmt.Sprintf("v = r * ω = 0.%d * 10 = %d m/s.", i%5+1, i%5+1),
				Options: []Option{
					{Text: fmt.Sprintf("%d m/s", i%5+1), IsCorrect: true},
					{Text: fmt.Sprintf("%d m/s", (i%5+1)*2), IsCorrect: false},
					{Text: fmt.Sprintf("%d m/s", (i%5+1)*10), IsCorrect: false},
					{Text: "0.5 m/s", IsCorrect: false},
				},
			}
		case 2:
			q = QuestionData{
				Text: fmt.Sprintf("Calculate the heat required to raise the temperature of %d kg of water by 10°C. (c = 4186 J/kg·K)", i%4+1),
				Topic: topic,
				Explanation: fmt.Sprintf("Q = m * c * ΔT = %d * 4186 * 10 = %d J.", i%4+1, (i%4+1)*41860),
				Options: []Option{
					{Text: fmt.Sprintf("%d J", (i%4+1)*41860), IsCorrect: true},
					{Text: fmt.Sprintf("%d J", (i%4+1)*4186), IsCorrect: false},
					{Text: "418600 J", IsCorrect: false},
					{Text: "83720 J", IsCorrect: false},
				},
			}
		case 3:
			q = QuestionData{
				Text: fmt.Sprintf("A spring with spring constant k = %d N/m is compressed by 0.2 m. What is the elastic potential energy stored?", (i%5+1)*100),
				Topic: topic,
				Explanation: fmt.Sprintf("PE = 0.5 * k * x² = 0.5 * %d * 0.04 = %.1f J.", (i%5+1)*100, float64((i%5+1)*100)*0.02),
				Options: []Option{
					{Text: fmt.Sprintf("%.1f J", float64((i%5+1)*100)*0.02), IsCorrect: true},
					{Text: fmt.Sprintf("%.1f J", float64((i%5+1)*100)*0.04), IsCorrect: false},
					{Text: "10.0 J", IsCorrect: false},
					{Text: "50.0 J", IsCorrect: false},
				},
			}
		case 4:
			q = QuestionData{
				Text: fmt.Sprintf("A convex lens has a focal length of %d cm. Where should an object be placed to form an image at the same size on the other side?", i%10+10),
				Topic: topic,
				Explanation: fmt.Sprintf("Object placed at 2f = 2 * %d = %d cm produces an inverted image of equal size at 2f.", i%10+10, 2*(i%10+10)),
				Options: []Option{
					{Text: fmt.Sprintf("%d cm (at 2f)", 2*(i%10+10)), IsCorrect: true},
					{Text: fmt.Sprintf("%d cm (at f)", i%10+10), IsCorrect: false},
					{Text: fmt.Sprintf("%d cm", 3*(i%10+10)), IsCorrect: false},
					{Text: "At infinity", IsCorrect: false},
				},
			}
		case 5:
			q = QuestionData{
				Text: fmt.Sprintf("A capacitor of %d μF is charged to a potential difference of 20 V. How much charge is stored?", i%5+2),
				Topic: topic,
				Explanation: fmt.Sprintf("Q = C * V = %d μF * 20 V = %d μC.", i%5+2, (i%5+2)*20),
				Options: []Option{
					{Text: fmt.Sprintf("%d μC", (i%5+2)*20), IsCorrect: true},
					{Text: fmt.Sprintf("%d μC", (i%5+2)*10), IsCorrect: false},
					{Text: fmt.Sprintf("%d μC", (i%5+2)*40), IsCorrect: false},
					{Text: "100 μC", IsCorrect: false},
				},
			}
		case 6:
			q = QuestionData{
				Text: fmt.Sprintf("A wire of length 2 m carrying a current of %d A is perpendicular to a uniform magnetic field of 0.5 T. What is the magnetic force?", i%4+1),
				Topic: topic,
				Explanation: fmt.Sprintf("F = I * L * B * sin(90°) = %d * 2 * 0.5 = %d N.", i%4+1, i%4+1),
				Options: []Option{
					{Text: fmt.Sprintf("%d N", i%4+1), IsCorrect: true},
					{Text: fmt.Sprintf("%d N", (i%4+1)*2), IsCorrect: false},
					{Text: fmt.Sprintf("%.1f N", float64(i%4+1)*0.5), IsCorrect: false},
					{Text: "10 N", IsCorrect: false},
				},
			}
		case 7:
			q = QuestionData{
				Text: fmt.Sprintf("What is the de Broglie wavelength of an electron moving with momentum %d × 10⁻²⁴ kg·m/s? (h = 6.63 × 10⁻³⁴ J·s)", i%5+2),
				Topic: topic,
				Explanation: fmt.Sprintf("λ = h / p = 6.63e-34 / (%de-24) = %.2f × 10⁻¹⁰ m.", i%5+2, 6.63/float64(i%5+2)),
				Options: []Option{
					{Text: fmt.Sprintf("%.2f × 10⁻¹⁰ m", 6.63/float64(i%5+2)), IsCorrect: true},
					{Text: fmt.Sprintf("%.2f × 10⁻¹² m", 6.63/float64(i%5+2)), IsCorrect: false},
					{Text: "3.31 × 10⁻¹⁰ m", IsCorrect: false},
					{Text: "1.00 × 10⁻⁹ m", IsCorrect: false},
				},
			}
		default:
			q = QuestionData{
				Text: fmt.Sprintf("In an ideal step-up transformer, the primary coil has 100 turns and the secondary has %d turns. If V_in = 120 V, what is V_out?", (i%4+2)*100),
				Topic: topic,
				Explanation: fmt.Sprintf("V_out / V_in = N_s / N_p => V_out = 120 * %d / 100 = %d V.", (i%4+2)*100, 120*(i%4+2)),
				Options: []Option{
					{Text: fmt.Sprintf("%d V", 120*(i%4+2)), IsCorrect: true},
					{Text: fmt.Sprintf("%d V", 60*(i%4+2)), IsCorrect: false},
					{Text: "120 V", IsCorrect: false},
					{Text: "2400 V", IsCorrect: false},
				},
			}
		}
		questions = append(questions, q)
	}

	return ExamData{
		Title:           "Grade 12 National Physics Mock Examination",
		Subject:         "Physics",
		Grade:           12,
		DurationMinutes: 120,
		PassMarks:       30,
		Instructions:    "This examination contains 60 multiple-choice questions covering all Grade 11-12 Physics units. Read every question carefully and select the best option.",
		Questions:       questions,
	}
}
