package main

func buildMathCourse() CourseData {
	return CourseData{
		Title:            "Grade 12 Mathematics: Ethiopian National University Entrance Prep",
		Slug:             "grade-12-mathematics-euee-prep",
		Subject:          "Mathematics",
		Grade:            12,
		Description:      "Master all key areas of Grade 11 & 12 Mathematics required for the Ethiopian University Entrance Examination (EUEE). Covers Vectors & 3D Geometry, Limits & Continuity, Differential & Integral Calculus, Sequences & Series, and Matrices with step-by-step problem-solving methods, unit quizzes, and national mock exams.",
		ShortDescription: "Complete Grade 12 National Exam Prep for Mathematics with vectors, calculus, series, quizzes, and mock exams.",
		ThumbnailUrl:     "https://images.unsplash.com/photo-1635070041078-e363dbe005cb?w=800&auto=format&fit=crop&q=80",
		Price:            0,
		IsFree:           true,
		Level:            "advanced",
		Sections: []SectionData{
			{
				Title:       "Unit 1: Sequences and Series",
				Description: "Arithmetic and Geometric sequences, limits of sequences, convergence tests, and financial applications.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 1.1: Arithmetic Progressions & Sum Formulas",
						Type:            "video",
						DurationSeconds: 1200,
						IsFreePreview:   true,
						Content: `# Arithmetic Progressions (AP)
An Arithmetic Progression is a sequence of numbers in which the difference between consecutive terms is constant.

### Key Formulas:
- **General Term (n-th term)**: $a_n = a_1 + (n - 1)d$
- **Sum of First n Terms**: $S_n = \frac{n}{2}[2a_1 + (n - 1)d] = \frac{n}{2}(a_1 + a_n)$
- **Common Difference**: $d = a_{k+1} - a_k$

### High-Yield Exam Tip:
When three terms are in AP, let them be $(a - d), a, (a + d)$. Their sum is $3a$, which instantly gives the middle term!`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 1.2: Geometric Progressions & Infinite Series",
						Type:            "note",
						DurationSeconds: 900,
						IsFreePreview:   true,
						Content: `# Geometric Progressions (GP) and Infinite Series
A sequence where the ratio of any two consecutive terms is constant.

### Key Formulas:
- **General Term**: $a_n = a_1 \cdot r^{n-1}$
- **Finite Sum**: $S_n = \frac{a_1(1 - r^n)}{1 - r}$ for $r \neq 1$
- **Sum to Infinity (Convergent GP)**: $S_\infty = \frac{a_1}{1 - r}$, valid only when $|r| < 1$.

### Critical Exam Rule:
If $|r| \ge 1$, the infinite geometric series diverges (no sum exists).`,
						PdfUrl: "sample.pdf",
					},
				},
			},
			{
				Title:       "Unit 2: Limits and Continuity",
				Description: "Intuitive and rigorous definitions of limits, one-sided limits, indeterminate forms (0/0), L'Hôpital's rule, and continuity on intervals.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 2.1: Evaluating Algebraic & Trigonometric Limits",
						Type:            "video",
						DurationSeconds: 1450,
						IsFreePreview:   true,
						Content: `# Evaluating Limits
Limits determine the behavior of functions near a point.

### Fundamental Limits:
1. $\lim_{x \to 0} \frac{\sin x}{x} = 1$
2. $\lim_{x \to 0} \frac{1 - \cos x}{x} = 0$
3. $\lim_{x \to \infty} \left(1 + \frac{1}{x}\right)^x = e$

### Indeterminate Forms & Factoring:
When encountering $\frac{0}{0}$, always factor polynomials or multiply by the conjugate radical before direct substitution.`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 2.2: Continuity, Intermediate Value Theorem & Asymptotes",
						Type:            "note",
						DurationSeconds: 850,
						IsFreePreview:   false,
						Content: `# Continuity and Asymptotes
A function $f(x)$ is continuous at $x = c$ if:
1. $f(c)$ is defined
2. $\lim_{x \to c} f(x)$ exists
3. $\lim_{x \to c} f(x) = f(c)$

### Types of Discontinuity:
- **Removable Discontinuity**: Hole in graph where limit exists but $\neq f(c)$
- **Jump Discontinuity**: Left and right limits are finite but unequal
- **Infinite Discontinuity**: Vertical asymptote exists where limit is $\pm\infty$`,
						PdfUrl: "sample.pdf",
					},
				},
			},
			{
				Title:       "Unit 3: Differential Calculus & Applications",
				Description: "Derivative rules, chain rule, implicit differentiation, curve sketching, related rates, and optimization.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 3.1: Derivatives Rules (Product, Quotient, Chain Rule)",
						Type:            "video",
						DurationSeconds: 1600,
						IsFreePreview:   false,
						Content: `# Differentiation Rules
- **Power Rule**: $\frac{d}{dx}[x^n] = n x^{n-1}$
- **Product Rule**: $\frac{d}{dx}[u \cdot v] = u'v + uv'$
- **Quotient Rule**: $\frac{d}{dx}\left[\frac{u}{v}\right] = \frac{u'v - uv'}{v^2}$
- **Chain Rule**: $\frac{d}{dx}[f(g(x))] = f'(g(x)) \cdot g'(x)$`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 3.2: Maxima, Minima, and Curve Sketching",
						Type:            "note",
						DurationSeconds: 1100,
						IsFreePreview:   false,
						Content: `# Extreme Values & Inflection Points
- **Critical Numbers**: Points where $f'(x) = 0$ or $f'(x)$ is undefined.
- **First Derivative Test**:
  - $f'$ changes from $+$ to $-$: Local Maximum
  - $f'$ changes from $-$ to $+$: Local Minimum
- **Concavity**:
  - $f''(x) > 0$: Concave Upward ($\cup$)
  - $f''(x) < 0$: Concave Downward ($\cap$)
- **Point of Inflection**: Where $f''(x) = 0$ and concavity changes.`,
						PdfUrl: "sample.pdf",
					},
				},
			},
			{
				Title:       "Unit 4: Vectors and 3D Analytic Geometry",
				Description: "Dot product, cross product, vector projections, equations of lines and planes in 3-space.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 4.1: Vector Operations, Dot Product & Orthogonality",
						Type:            "video",
						DurationSeconds: 1350,
						IsFreePreview:   false,
						Content: `# Vectors in 3-Dimensional Space
Given $\vec{u} = (u_1, u_2, u_3)$ and $\vec{v} = (v_1, v_2, v_3)$:

### Dot Product:
$\vec{u} \cdot \vec{v} = u_1 v_1 + u_2 v_2 + u_3 v_3 = \|\vec{u}\| \|\vec{v}\| \cos\theta$

### Orthogonality Condition:
Two non-zero vectors are perpendicular ($\theta = 90^\circ$) if and only if:
$$\vec{u} \cdot \vec{v} = 0$$`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 4.2: Cross Product & Equations of Lines and Planes",
						Type:            "note",
						DurationSeconds: 950,
						IsFreePreview:   false,
						Content: `# Cross Product & 3D Planes
The cross product $\vec{u} \times \vec{v}$ produces a vector perpendicular to both $\vec{u}$ and $\vec{v}$.

### Equation of a Plane:
Given normal vector $\vec{n} = (A, B, C)$ and point $P_0(x_0, y_0, z_0)$:
$$A(x - x_0) + B(y - y_0) + C(z - z_0) = 0$$
Standard Form: $Ax + By + Cz + D = 0$.`,
						PdfUrl: "sample.pdf",
					},
				},
			},
		},
		Quizzes: []QuizData{
			{
				Title:            "Math Quiz 1: Sequences & Series Mastery",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "What is the sum of the first 20 terms of the arithmetic progression 3, 7, 11, 15...?",
						Explanation: "Here a1 = 3, d = 4, n = 20. S20 = (20/2)[2(3) + (20-1)(4)] = 10[6 + 76] = 10(82) = 820.",
						Options: []Option{
							{Text: "820", IsCorrect: true},
							{Text: "780", IsCorrect: false},
							{Text: "840", IsCorrect: false},
							{Text: "800", IsCorrect: false},
						},
					},
					{
						Text:        "For which range of values of x does the infinite series 1 + (x-2) + (x-2)^2 + ... converge?",
						Explanation: "A geometric series converges when |r| < 1. Here r = (x-2). So |x - 2| < 1 => -1 < x - 2 < 1 => 1 < x < 3.",
						Options: []Option{
							{Text: "1 < x < 3", IsCorrect: true},
							{Text: "-1 < x < 1", IsCorrect: false},
							{Text: "0 < x < 2", IsCorrect: false},
							{Text: "x > 2", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Math Quiz 2: Limits & Continuity",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "Evaluate lim(x->3) (x^2 - 9) / (x - 3).",
						Explanation: "Factor numerator: (x - 3)(x + 3)/(x - 3) = x + 3. As x -> 3, limit is 3 + 3 = 6.",
						Options: []Option{
							{Text: "6", IsCorrect: true},
							{Text: "0", IsCorrect: false},
							{Text: "3", IsCorrect: false},
							{Text: "Undefined", IsCorrect: false},
						},
					},
					{
						Text:        "What is the value of lim(x->0) (sin 5x) / (2x)?",
						Explanation: "lim(x->0) (sin 5x)/(2x) = (5/2) * lim(5x->0) (sin 5x)/(5x) = 5/2 * 1 = 5/2 = 2.5.",
						Options: []Option{
							{Text: "5/2", IsCorrect: true},
							{Text: "1", IsCorrect: false},
							{Text: "2/5", IsCorrect: false},
							{Text: "5", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Math Quiz 3: Differential Calculus Rules",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "If f(x) = x^3 * e^(2x), what is f'(x)?",
						Explanation: "Use product rule: f'(x) = 3x^2 * e^(2x) + x^3 * 2e^(2x) = e^(2x)(3x^2 + 2x^3) = x^2 e^(2x)(3 + 2x).",
						Options: []Option{
							{Text: "x^2 e^(2x)(3 + 2x)", IsCorrect: true},
							{Text: "3x^2 e^(2x)", IsCorrect: false},
							{Text: "2x^3 e^(2x)", IsCorrect: false},
							{Text: "6x^2 e^(2x)", IsCorrect: false},
						},
					},
					{
						Text:        "Find the slope of the tangent line to the curve y = ln(x^2 + 1) at x = 1.",
						Explanation: "dy/dx = 2x / (x^2 + 1). At x = 1, dy/dx = 2(1)/(1+1) = 2/2 = 1.",
						Options: []Option{
							{Text: "1", IsCorrect: true},
							{Text: "2", IsCorrect: false},
							{Text: "1/2", IsCorrect: false},
							{Text: "0", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Math Quiz 4: Curve Sketching & Optimization",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "At what value of x does f(x) = 2x^3 - 9x^2 + 12x + 5 have a local minimum?",
						Explanation: "f'(x) = 6x^2 - 18x + 12 = 6(x^2 - 3x + 2) = 6(x - 1)(x - 2) = 0 => critical points at x = 1, 2. f''(x) = 12x - 18. At x = 2, f''(2) = 24 - 18 = 6 > 0 (Local Min).",
						Options: []Option{
							{Text: "x = 2", IsCorrect: true},
							{Text: "x = 1", IsCorrect: false},
							{Text: "x = 3", IsCorrect: false},
							{Text: "x = 0", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Math Quiz 5: Vectors and 3D Geometry",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "If vector u = (2, -1, 3) and vector v = (k, 4, 2) are perpendicular, find the value of k.",
						Explanation: "u . v = 0 => 2(k) + (-1)(4) + 3(2) = 0 => 2k - 4 + 6 = 0 => 2k + 2 = 0 => k = -1.",
						Options: []Option{
							{Text: "-1", IsCorrect: true},
							{Text: "1", IsCorrect: false},
							{Text: "-2", IsCorrect: false},
							{Text: "2", IsCorrect: false},
						},
					},
				},
			},
		},
		Exams: []ExamData{
			{
				Title:           "Mathematics Unit 1 & 2 Progress Examination",
				Subject:         "Mathematics",
				Grade:           12,
				DurationMinutes: 45,
				PassMarks:       12,
				Instructions:    "Covers Sequences & Series, Limits, and Continuity. Calculators permitted. 45 minutes total.",
				Questions: []QuestionData{
					{
						Text:        "Find the 15th term of the sequence: 4, 9, 14, 19...",
						Explanation: "a15 = a1 + 14d = 4 + 14(5) = 4 + 70 = 74.",
						Options: []Option{
							{Text: "74", IsCorrect: true},
							{Text: "79", IsCorrect: false},
							{Text: "69", IsCorrect: false},
							{Text: "84", IsCorrect: false},
						},
					},
					{
						Text:        "What is the sum of the infinite series 6 + 2 + 2/3 + 2/9 + ...?",
						Explanation: "a = 6, r = 1/3. S = 6 / (1 - 1/3) = 6 / (2/3) = 9.",
						Options: []Option{
							{Text: "9", IsCorrect: true},
							{Text: "8", IsCorrect: false},
							{Text: "12", IsCorrect: false},
							{Text: "18", IsCorrect: false},
						},
					},
					{
						Text:        "Evaluate: lim(x->4) (sqrt(x) - 2) / (x - 4).",
						Explanation: "Multiply by conjugate: (sqrt(x)-2)(sqrt(x)+2)/((x-4)(sqrt(x)+2)) = (x-4)/((x-4)(sqrt(x)+2)) = 1/(sqrt(x)+2). At x=4, 1/(2+2) = 1/4.",
						Options: []Option{
							{Text: "1/4", IsCorrect: true},
							{Text: "1/2", IsCorrect: false},
							{Text: "0", IsCorrect: false},
							{Text: "1/8", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:           "Mathematics Midterm Examination (Calculus & Vectors)",
				Subject:         "Mathematics",
				Grade:           12,
				DurationMinutes: 60,
				PassMarks:       15,
				Instructions:    "Midterm Exam testing Differential Calculus, Optimization, and 3D Vectors.",
				Questions: []QuestionData{
					{
						Text:        "If y = e^(sin x), what is dy/dx at x = 0?",
						Explanation: "dy/dx = e^(sin x) * cos x. At x = 0, e^0 * cos 0 = 1 * 1 = 1.",
						Options: []Option{
							{Text: "1", IsCorrect: true},
							{Text: "0", IsCorrect: false},
							{Text: "e", IsCorrect: false},
							{Text: "-1", IsCorrect: false},
						},
					},
					{
						Text:        "Find the magnitude of the vector projection of u = (3, 4) onto v = (1, 0).",
						Explanation: "proj = |u . v| / ||v|| = |3(1) + 4(0)| / 1 = 3.",
						Options: []Option{
							{Text: "3", IsCorrect: true},
							{Text: "4", IsCorrect: false},
							{Text: "5", IsCorrect: false},
							{Text: "1", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:           "Ethiopian National EUEE Mathematics Final Mock Exam",
				Subject:         "Mathematics",
				Grade:           12,
				DurationMinutes: 90,
				PassMarks:       25,
				Instructions:    "Comprehensive National Entrance Examination simulation for Grade 12 Natural & Social Science candidates.",
				Questions: []QuestionData{
					{
						Text:        "What is the derivative of f(x) = (2x + 1)^4 at x = 0?",
						Explanation: "f'(x) = 4(2x + 1)^3 * 2 = 8(2x + 1)^3. At x = 0, 8(1)^3 = 8.",
						Options: []Option{
							{Text: "8", IsCorrect: true},
							{Text: "4", IsCorrect: false},
							{Text: "16", IsCorrect: false},
							{Text: "2", IsCorrect: false},
						},
					},
					{
						Text:        "Find the distance from the point (1, 2, 3) to the plane 2x - y + 2z - 4 = 0.",
						Explanation: "d = |2(1) - 1(2) + 2(3) - 4| / sqrt(2^2 + (-1)^2 + 2^2) = |2 - 2 + 6 - 4| / sqrt(4 + 1 + 4) = 2 / sqrt(9) = 2/3.",
						Options: []Option{
							{Text: "2/3", IsCorrect: true},
							{Text: "1", IsCorrect: false},
							{Text: "4/3", IsCorrect: false},
							{Text: "2", IsCorrect: false},
						},
					},
				},
			},
		},
	}
}
