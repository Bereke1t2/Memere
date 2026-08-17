package main

func buildPhysicsCourse() CourseData {
	return CourseData{
		Title:            "Grade 12 Physics: Ethiopian National University Entrance Prep",
		Slug:             "grade-12-physics-euee-prep",
		Subject:          "Physics",
		Grade:            12,
		Description:      "Comprehensive Grade 11 & 12 Physics course engineered for the Ethiopian University Entrance Examination. Deep coverage of 2D Kinematics & Vectors, Rotational Dynamics, Electromagnetism & Induction, Optics, Wave Motion, and Atomic/Nuclear Physics with past national exam questions, unit quizzes, and mock tests.",
		ShortDescription: "Grade 12 Physics EUEE prep covering mechanics, electromagnetism, optics, modern physics, quizzes, and mock exams.",
		ThumbnailUrl:     "https://images.unsplash.com/photo-1507668077129-56e32842fceb?w=800&auto=format&fit=crop&q=80",
		Price:            0,
		IsFree:           true,
		Level:            "advanced",
		Sections: []SectionData{
			{
				Title:       "Unit 1: Two-Dimensional Motion and Vectors",
				Description: "Projectile motion equations, maximum height, range, flight time, circular motion, and relative velocity.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 1.1: Projectile Motion & Trajectory Equations",
						Type:            "video",
						DurationSeconds: 1300,
						IsFreePreview:   true,
						Content: `# Projectile Motion in 2 Dimensions
A projectile moves under constant downward gravitational acceleration ($g = 9.8\text{ m/s}^2$).

### Horizontal Motion:
- $a_x = 0$
- $v_x = v_0 \cos\theta = \text{constant}$
- $x(t) = (v_0 \cos\theta) t$

### Vertical Motion:
- $a_y = -g$
- $v_y(t) = v_0 \sin\theta - gt$
- $y(t) = (v_0 \sin\theta) t - \frac{1}{2}gt^2$

### Key Exam Formulas:
- **Time of Flight**: $T = \frac{2v_0 \sin\theta}{g}$
- **Maximum Height**: $H_{\max} = \frac{(v_0 \sin\theta)^2}{2g}$
- **Horizontal Range**: $R = \frac{v_0^2 \sin(2\theta)}{g}$ (Maximum at $\theta = 45^\circ$)`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 1.2: Uniform Circular Motion & Centripetal Dynamics",
						Type:            "note",
						DurationSeconds: 850,
						IsFreePreview:   true,
						Content: `# Uniform Circular Motion
An object traveling in a circle of radius $r$ at constant speed $v$ experiences a continuous change in direction.

### Centripetal Acceleration:
$$a_c = \frac{v^2}{r} = \omega^2 r$$

### Centripetal Force:
$$F_c = m a_c = \frac{m v^2}{r}$$

### Banked Curves:
For an unbanked road with static friction: $v_{\max} = \sqrt{\mu_s g r}$.
For a frictionless banked curve at angle $\theta$: $\tan\theta = \frac{v^2}{rg}$.`,
						PdfUrl: "sample.pdf",
					},
				},
			},
			{
				Title:       "Unit 2: Rotational Dynamics & Equilibrium",
				Description: "Torque, moment of inertia, angular momentum conservation, and conditions for static equilibrium.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 2.1: Torque & Static Equilibrium Conditions",
						Type:            "video",
						DurationSeconds: 1400,
						IsFreePreview:   false,
						Content: `# Torque and Static Equilibrium
Torque measures the rotational effect of a force applied at distance $r$ from an axis:
$$\tau = \vec{r} \times \vec{F} = r F \sin\theta$$

### Conditions for Static Equilibrium:
1. **Translational Equilibrium**: $\sum \vec{F} = 0$ ($\sum F_x = 0$, $\sum F_y = 0$)
2. **Rotational Equilibrium**: $\sum \vec{\tau} = 0$ about any pivot point.`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 2.2: Moment of Inertia & Angular Momentum",
						Type:            "note",
						DurationSeconds: 950,
						IsFreePreview:   false,
						Content: `# Angular Momentum Conservation
- **Newton's Second Law for Rotation**: $\tau_{\text{net}} = I \alpha$
- **Rotational Kinetic Energy**: $K_{\text{rot}} = \frac{1}{2} I \omega^2$
- **Angular Momentum**: $L = I \omega$

### Conservation Law:
When external net torque $\tau_{\text{ext}} = 0$:
$$I_1 \omega_1 = I_2 \omega_2 = \text{constant}$$`,
						PdfUrl: "sample.pdf",
					},
				},
			},
			{
				Title:       "Unit 3: Electromagnetism & Electromagnetic Induction",
				Description: "Coulomb's Law, electric potential, magnetic fields, Lorentz force, Faraday's Law, and Lenz's Law.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 3.1: Magnetic Fields & Lorentz Force",
						Type:            "video",
						DurationSeconds: 1550,
						IsFreePreview:   false,
						Content: `# Magnetic Force on Charges & Currents
- **Force on a Moving Charge**: $\vec{F}_B = q(\vec{v} \times \vec{B}) \implies F = q v B \sin\theta$
- **Force on a Current-Carrying Wire**: $\vec{F} = I(\vec{L} \times \vec{B}) \implies F = I L B \sin\theta$
- **Radius of Circular Path in Uniform Magnetic Field**:
  $$q v B = \frac{m v^2}{r} \implies r = \frac{m v}{q B}$$`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 3.2: Faraday's Law, Lenz's Law & Transformers",
						Type:            "note",
						DurationSeconds: 1000,
						IsFreePreview:   false,
						Content: `# Faraday's Law of Electromagnetic Induction
- **Magnetic Flux**: $\Phi_B = \vec{B} \cdot \vec{A} = B A \cos\theta$
- **Induced EMF**: $\mathcal{E} = -N \frac{d\Phi_B}{dt}$
- **Lenz's Law**: The direction of any induced current opposes the change in flux producing it.

### Transformer Equation:
$$\frac{V_s}{V_p} = \frac{N_s}{N_p} = \frac{I_p}{I_s}$$`,
						PdfUrl: "sample.pdf",
					},
				},
			},
			{
				Title:       "Unit 4: Modern Physics & Nuclear Reactions",
				Description: "Photoelectric effect, de Broglie wavelength, Bohr model of atom, radioactive decay, and mass-energy equivalence.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 4.1: The Photoelectric Effect & Photon Energy",
						Type:            "video",
						DurationSeconds: 1250,
						IsFreePreview:   false,
						Content: `# Photoelectric Effect
Einstein's photoelectric equation relates photon energy to the metal work function ($\Phi$) and the maximum kinetic energy of emitted electrons:
$$E = h f = \Phi + K_{\max} = h f_0 + e V_0$$

### Key Exam Insights:
1. Increasing light intensity increases the *number* of photoelectrons (photocurrent), but NOT their kinetic energy.
2. Increasing light frequency increases the *kinetic energy* of photoelectrons.`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 4.2: Radioactive Decay Law & Half-Life Calculations",
						Type:            "note",
						DurationSeconds: 900,
						IsFreePreview:   false,
						Content: `# Nuclear Physics & Radioactive Decay
- **Decay Law**: $N(t) = N_0 e^{-\lambda t} = N_0 \left(\frac{1}{2}\right)^{t / T_{1/2}}$
- **Half-Life Formula**: $T_{1/2} = \frac{\ln 2}{\lambda} \approx \frac{0.693}{\lambda}$
- **Mass Defect & Binding Energy**: $E = (\Delta m) c^2$`,
						PdfUrl: "sample.pdf",
					},
				},
			},
		},
		Quizzes: []QuizData{
			{
				Title:            "Physics Quiz 1: 2D Kinematics & Projectiles",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "At what angle of launch is the horizontal range of an ideal projectile maximized?",
						Explanation: "Range formula R = (v0^2 * sin 2theta)/g. sin(2theta) is maximum when 2theta = 90 deg, so theta = 45 deg.",
						Options: []Option{
							{Text: "45 degrees", IsCorrect: true},
							{Text: "30 degrees", IsCorrect: false},
							{Text: "60 degrees", IsCorrect: false},
							{Text: "90 degrees", IsCorrect: false},
						},
					},
					{
						Text:        "A ball is thrown horizontally with 10 m/s from a cliff. How long does it take to fall 45 m vertically? (g = 10 m/s^2)",
						Explanation: "y = 1/2 g t^2 => 45 = 1/2(10) t^2 = 5 t^2 => t^2 = 9 => t = 3 seconds.",
						Options: []Option{
							{Text: "3.0 s", IsCorrect: true},
							{Text: "4.5 s", IsCorrect: false},
							{Text: "2.0 s", IsCorrect: false},
							{Text: "9.0 s", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Physics Quiz 2: Circular Motion & Gravitation",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "If the speed of an object moving in a circle of constant radius is doubled, what happens to the centripetal force?",
						Explanation: "Fc = m*v^2 / r. Since Fc is proportional to v^2, doubling v increases Fc by a factor of 4.",
						Options: []Option{
							{Text: "It quadruples", IsCorrect: true},
							{Text: "It doubles", IsCorrect: false},
							{Text: "It halves", IsCorrect: false},
							{Text: "It remains unchanged", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Physics Quiz 3: Rotational Equilibrium & Torque",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "A uniform meter rule of weight 2 N is pivoted at the 50 cm mark. A 4 N weight is hung at 20 cm. Where must a 6 N weight be placed to balance it?",
						Explanation: "Clockwise torque = Counterclockwise torque. 4 N * (50 - 20) cm = 6 N * (x - 50) cm => 4 * 30 = 6 * (x - 50) => 120 = 6(x - 50) => 20 = x - 50 => x = 70 cm mark.",
						Options: []Option{
							{Text: "At the 70 cm mark", IsCorrect: true},
							{Text: "At the 80 cm mark", IsCorrect: false},
							{Text: "At the 60 cm mark", IsCorrect: false},
							{Text: "At the 90 cm mark", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Physics Quiz 4: Electromagnetism & Induction",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "An electron enters a magnetic field directed perpendicularly into the page. What is the shape of its trajectory in the field?",
						Explanation: "The magnetic Lorentz force acts perpendicular to the velocity at all times, providing centripetal acceleration, resulting in a circular arc.",
						Options: []Option{
							{Text: "A circle or circular arc", IsCorrect: true},
							{Text: "A straight line", IsCorrect: false},
							{Text: "A parabola", IsCorrect: false},
							{Text: "A hyperbola", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Physics Quiz 5: Photoelectric Effect & Modern Physics",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "In a photoelectric experiment, doubling the intensity of incident light above threshold frequency causes:",
						Explanation: "Light intensity is proportional to photon count. Doubling intensity doubles the number of emitted photoelectrons per second, leaving individual kinetic energies unchanged.",
						Options: []Option{
							{Text: "Double the number of emitted electrons", IsCorrect: true},
							{Text: "Double the kinetic energy of emitted electrons", IsCorrect: false},
							{Text: "Double the stopping potential", IsCorrect: false},
							{Text: "Double the work function", IsCorrect: false},
						},
					},
				},
			},
		},
		Exams: []ExamData{
			{
				Title:           "Physics Unit 1 & 2 Mechanics Examination",
				Subject:         "Physics",
				Grade:           12,
				DurationMinutes: 45,
				PassMarks:       12,
				Instructions:    "Covers 2D Kinematics, Circular Motion, Torque, and Rotational Dynamics. 45 minutes duration.",
				Questions: []QuestionData{
					{
						Text:        "A projectile is fired at 30 degrees with an initial velocity of 40 m/s. What is its vertical velocity at the highest point?",
						Explanation: "At maximum height, the vertical velocity component vy is always zero.",
						Options: []Option{
							{Text: "0 m/s", IsCorrect: true},
							{Text: "20 m/s", IsCorrect: false},
							{Text: "34.6 m/s", IsCorrect: false},
							{Text: "40 m/s", IsCorrect: false},
						},
					},
					{
						Text:        "What is the angular momentum of a 2 kg mass rotating with angular velocity 4 rad/s at radius 0.5 m?",
						Explanation: "I = m*r^2 = 2 * (0.5)^2 = 2 * 0.25 = 0.5 kg m^2. L = I * omega = 0.5 * 4 = 2.0 kg m^2/s.",
						Options: []Option{
							{Text: "2.0 kg m^2/s", IsCorrect: true},
							{Text: "4.0 kg m^2/s", IsCorrect: false},
							{Text: "1.0 kg m^2/s", IsCorrect: false},
							{Text: "8.0 kg m^2/s", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:           "Physics Midterm Examination (Electromagnetism & Circuits)",
				Subject:         "Physics",
				Grade:           12,
				DurationMinutes: 60,
				PassMarks:       15,
				Instructions:    "Testing Magnetic Fields, Lorentz Force, Induction, and AC Circuits.",
				Questions: []QuestionData{
					{
						Text:        "A step-down transformer has a primary coil of 1000 turns and secondary coil of 200 turns. If the input primary voltage is 220 V, what is the secondary output voltage?",
						Explanation: "Vs / Vp = Ns / Np => Vs / 220 = 200 / 1000 => Vs = 220 * 0.2 = 44 V.",
						Options: []Option{
							{Text: "44 V", IsCorrect: true},
							{Text: "110 V", IsCorrect: false},
							{Text: "22 V", IsCorrect: false},
							{Text: "1100 V", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:           "Ethiopian National EUEE Physics Final Mock Exam",
				Subject:         "Physics",
				Grade:           12,
				DurationMinutes: 90,
				PassMarks:       25,
				Instructions:    "Comprehensive National Entrance Examination simulation for Grade 12 Natural Science Physics.",
				Questions: []QuestionData{
					{
						Text:        "A radioactive substance has a half-life of 8 days. If the initial sample mass is 80 g, how much remains undecayed after 24 days?",
						Explanation: "Number of half-lives n = 24 / 8 = 3. Remaining mass = 80 / (2^3) = 80 / 8 = 10 g.",
						Options: []Option{
							{Text: "10 g", IsCorrect: true},
							{Text: "20 g", IsCorrect: false},
							{Text: "5 g", IsCorrect: false},
							{Text: "2.5 g", IsCorrect: false},
						},
					},
				},
			},
		},
	}
}
