package main

import "fmt"

func buildMathExam() ExamData {
	questions := []QuestionData{
		{
			Text: "What is the domain of the function f(x) = √(4 - x²)?",
			Topic: "Functions",
			Explanation: "For the square root to be real, 4 - x² >= 0, which means x² <= 4, giving [-2, 2].",
			Options: []Option{
				{Text: "[-2, 2]", IsCorrect: true},
				{Text: "(-∞, -2] U [2, ∞)", IsCorrect: false},
				{Text: "[0, 2]", IsCorrect: false},
				{Text: "(-2, 2)", IsCorrect: false},
			},
		},
		{
			Text: "If f(x) = 2x + 3 and g(x) = x² - 1, what is (f ∘ g)(2)?",
			Topic: "Functions",
			Explanation: "g(2) = 2² - 1 = 3. Then f(g(2)) = f(3) = 2(3) + 3 = 9.",
			Options: []Option{
				{Text: "9", IsCorrect: true},
				{Text: "7", IsCorrect: false},
				{Text: "11", IsCorrect: false},
				{Text: "6", IsCorrect: false},
			},
		},
		{
			Text: "Which of the following is the inverse function f⁻¹(x) of f(x) = (3x - 5) / 2?",
			Topic: "Functions",
			Explanation: "Let y = (3x - 5)/2 => 2y = 3x - 5 => x = (2y + 5)/3. Thus f⁻¹(x) = (2x + 5)/3.",
			Options: []Option{
				{Text: "f⁻¹(x) = (2x + 5) / 3", IsCorrect: true},
				{Text: "f⁻¹(x) = (2x - 5) / 3", IsCorrect: false},
				{Text: "f⁻¹(x) = (3x + 5) / 2", IsCorrect: false},
				{Text: "f⁻¹(x) = 2 / (3x - 5)", IsCorrect: false},
			},
		},
		{
			Text: "What is the sum of the first 20 terms of the arithmetic progression 3, 7, 11, 15, ...?",
			Topic: "Sequences and Series",
			Explanation: "a = 3, d = 4. S₂₀ = (20/2)[2(3) + (20-1)4] = 10[6 + 76] = 820.",
			Options: []Option{
				{Text: "820", IsCorrect: true},
				{Text: "800", IsCorrect: false},
				{Text: "840", IsCorrect: false},
				{Text: "780", IsCorrect: false},
			},
		},
		{
			Text: "What is the sum of the infinite geometric series 12 + 6 + 3 + 1.5 + ...?",
			Topic: "Sequences and Series",
			Explanation: "a = 12, r = 0.5. S_inf = a / (1 - r) = 12 / (1 - 0.5) = 24.",
			Options: []Option{
				{Text: "24", IsCorrect: true},
				{Text: "20", IsCorrect: false},
				{Text: "36", IsCorrect: false},
				{Text: "18", IsCorrect: false},
			},
		},
		{
			Text: "What is the limit as x approaches 2 of (x² - 4) / (x - 2)?",
			Topic: "Limits and Continuity",
			Explanation: "(x² - 4)/(x - 2) = (x - 2)(x + 2)/(x - 2) = x + 2. As x -> 2, limit = 4.",
			Options: []Option{
				{Text: "4", IsCorrect: true},
				{Text: "2", IsCorrect: false},
				{Text: "0", IsCorrect: false},
				{Text: "Undefined", IsCorrect: false},
			},
		},
		{
			Text: "What is the limit as x approaches 0 of sin(5x) / x?",
			Topic: "Limits and Continuity",
			Explanation: "lim_{x->0} sin(kx)/x = k. Here k = 5, so the limit is 5.",
			Options: []Option{
				{Text: "5", IsCorrect: true},
				{Text: "1", IsCorrect: false},
				{Text: "0", IsCorrect: false},
				{Text: "1/5", IsCorrect: false},
			},
		},
		{
			Text: "What is the derivative of f(x) = 3x⁴ - 5x² + 7x - 2 with respect to x?",
			Topic: "Derivatives",
			Explanation: "f'(x) = 12x³ - 10x + 7.",
			Options: []Option{
				{Text: "12x³ - 10x + 7", IsCorrect: true},
				{Text: "12x³ - 10x² + 7", IsCorrect: false},
				{Text: "7x³ - 10x + 7", IsCorrect: false},
				{Text: "12x⁴ - 10x² + 7", IsCorrect: false},
			},
		},
		{
			Text: "If y = sin(x²), what is dy/dx?",
			Topic: "Derivatives",
			Explanation: "Using chain rule: dy/dx = cos(x²) * d(x²)/dx = 2x cos(x²).",
			Options: []Option{
				{Text: "2x cos(x²)", IsCorrect: true},
				{Text: "cos(x²)", IsCorrect: false},
				{Text: "2x sin(x²)", IsCorrect: false},
				{Text: "-2x cos(x²)", IsCorrect: false},
			},
		},
		{
			Text: "What is the slope of the tangent line to the curve y = x³ - 3x + 2 at x = 2?",
			Topic: "Derivatives",
			Explanation: "dy/dx = 3x² - 3. At x = 2: 3(4) - 3 = 9.",
			Options: []Option{
				{Text: "9", IsCorrect: true},
				{Text: "6", IsCorrect: false},
				{Text: "12", IsCorrect: false},
				{Text: "3", IsCorrect: false},
			},
		},
		{
			Text: "Evaluate the definite integral ∫ from 0 to 3 of (2x + 1) dx.",
			Topic: "Integrals",
			Explanation: "∫(2x + 1)dx = x² + x evaluated from 0 to 3: (3² + 3) - 0 = 9 + 3 = 12.",
			Options: []Option{
				{Text: "12", IsCorrect: true},
				{Text: "9", IsCorrect: false},
				{Text: "15", IsCorrect: false},
				{Text: "6", IsCorrect: false},
			},
		},
		{
			Text: "What is the indefinite integral ∫ e^(3x) dx?",
			Topic: "Integrals",
			Explanation: "∫ e^(kx) dx = (1/k)e^(kx) + C. Here (1/3)e^(3x) + C.",
			Options: []Option{
				{Text: "(1/3)e^(3x) + C", IsCorrect: true},
				{Text: "3e^(3x) + C", IsCorrect: false},
				{Text: "e^(3x) + C", IsCorrect: false},
				{Text: "(1/3)e^x + C", IsCorrect: false},
			},
		},
		{
			Text: "What is the distance between the points P(1, 2, 3) and Q(4, 6, 3) in 3D space?",
			Topic: "Vectors and 3D Geometry",
			Explanation: "d = √((4-1)² + (6-2)² + (3-3)²) = √(9 + 16 + 0) = √25 = 5.",
			Options: []Option{
				{Text: "5", IsCorrect: true},
				{Text: "7", IsCorrect: false},
				{Text: "√29", IsCorrect: false},
				{Text: "4", IsCorrect: false},
			},
		},
		{
			Text: "If vector u = (2, 3) and vector v = (4, -1), what is the dot product u · v?",
			Topic: "Vectors and 3D Geometry",
			Explanation: "u · v = 2(4) + 3(-1) = 8 - 3 = 5.",
			Options: []Option{
				{Text: "5", IsCorrect: true},
				{Text: "11", IsCorrect: false},
				{Text: "-5", IsCorrect: false},
				{Text: "8", IsCorrect: false},
			},
		},
		{
			Text: "What is the determinant of the matrix A = [[3, 2], [1, 4]]?",
			Topic: "Matrices",
			Explanation: "det(A) = 3(4) - 2(1) = 12 - 2 = 10.",
			Options: []Option{
				{Text: "10", IsCorrect: true},
				{Text: "14", IsCorrect: false},
				{Text: "12", IsCorrect: false},
				{Text: "-10", IsCorrect: false},
			},
		},
		{
			Text: "If A is a 2x2 matrix with det(A) = 5, what is det(2A)?",
			Topic: "Matrices",
			Explanation: "For an n x n matrix, det(kA) = kⁿ det(A). For n = 2, det(2A) = 2² * 5 = 20.",
			Options: []Option{
				{Text: "20", IsCorrect: true},
				{Text: "10", IsCorrect: false},
				{Text: "40", IsCorrect: false},
				{Text: "5", IsCorrect: false},
			},
		},
		{
			Text: "A bag contains 4 red balls and 6 blue balls. What is the probability of drawing a red ball?",
			Topic: "Probability",
			Explanation: "P(Red) = 4 / (4 + 6) = 4/10 = 2/5 = 0.4.",
			Options: []Option{
				{Text: "2/5", IsCorrect: true},
				{Text: "3/5", IsCorrect: false},
				{Text: "1/4", IsCorrect: false},
				{Text: "1/2", IsCorrect: false},
			},
		},
		{
			Text: "What is the value of ₇P₃ (permutations of 7 items taken 3 at a time)?",
			Topic: "Probability",
			Explanation: "7! / (7 - 3)! = 7 * 6 * 5 = 210.",
			Options: []Option{
				{Text: "210", IsCorrect: true},
				{Text: "35", IsCorrect: false},
				{Text: "140", IsCorrect: false},
				{Text: "42", IsCorrect: false},
			},
		},
		{
			Text: "What is the mean of the data set: 4, 8, 6, 5, 3, 4?",
			Topic: "Statistics",
			Explanation: "Sum = 30, Count = 6. Mean = 30 / 6 = 5.",
			Options: []Option{
				{Text: "5", IsCorrect: true},
				{Text: "4.5", IsCorrect: false},
				{Text: "6", IsCorrect: false},
				{Text: "4", IsCorrect: false},
			},
		},
		{
			Text: "What is the median of the data set: 12, 7, 3, 9, 15, 8, 4?",
			Topic: "Statistics",
			Explanation: "Ordered: 3, 4, 7, 8, 9, 12, 15. The middle value (4th item) is 8.",
			Options: []Option{
				{Text: "8", IsCorrect: true},
				{Text: "7", IsCorrect: false},
				{Text: "9", IsCorrect: false},
				{Text: "8.5", IsCorrect: false},
			},
		},
	}

	// Generate remaining 40 math questions programmatically across topics
	mathTopics := []string{
		"Algebra & Polynomials", "Trigonometry", "Coordinate Geometry", "Calculus & Optimization",
		"Complex Numbers", "Vectors & Matrices", "Probability & Distributions", "Analytical Geometry",
	}

	for i := 21; i <= 60; i++ {
		topic := mathTopics[i%len(mathTopics)]
		var q QuestionData
		switch i % 8 {
		case 1:
			q = QuestionData{
				Text: fmt.Sprintf("Solve for x: log₂(x + %d) - log₂(x - 1) = 1.", i),
				Topic: topic,
				Explanation: fmt.Sprintf("(x + %d)/(x - 1) = 2 => x + %d = 2x - 2 => x = %d.", i, i, i+2),
				Options: []Option{
					{Text: fmt.Sprintf("%d", i+2), IsCorrect: true},
					{Text: fmt.Sprintf("%d", i+1), IsCorrect: false},
					{Text: fmt.Sprintf("%d", i-1), IsCorrect: false},
					{Text: fmt.Sprintf("%d", 2*i), IsCorrect: false},
				},
			}
		case 2:
			q = QuestionData{
				Text: fmt.Sprintf("What is the exact value of sin(%d° + 30°)?", (i*15)%90),
				Topic: topic,
				Explanation: "Using compound angle identities: sin(A + B) = sin A cos B + cos A sin B.",
				Options: []Option{
					{Text: "(√6 + √2) / 4", IsCorrect: true},
					{Text: "(√6 - √2) / 4", IsCorrect: false},
					{Text: "1/2", IsCorrect: false},
					{Text: "√3 / 2", IsCorrect: false},
				},
			}
		case 3:
			q = QuestionData{
				Text: fmt.Sprintf("Find the equation of the line perpendicular to 2x + 3y = %d passing through (1, 2).", i),
				Topic: topic,
				Explanation: "Slope of original line is -2/3. Perpendicular slope is 3/2. y - 2 = (3/2)(x - 1) => 3x - 2y + 1 = 0.",
				Options: []Option{
					{Text: "3x - 2y + 1 = 0", IsCorrect: true},
					{Text: "2x + 3y - 8 = 0", IsCorrect: false},
					{Text: "3x + 2y - 7 = 0", IsCorrect: false},
					{Text: "2x - 3y + 4 = 0", IsCorrect: false},
				},
			}
		case 4:
			q = QuestionData{
				Text: fmt.Sprintf("Find the critical point of f(x) = x² - %dx + 10.", 2*i),
				Topic: topic,
				Explanation: fmt.Sprintf("f'(x) = 2x - %d = 0 => x = %d.", 2*i, i),
				Options: []Option{
					{Text: fmt.Sprintf("x = %d", i), IsCorrect: true},
					{Text: fmt.Sprintf("x = %d", 2*i), IsCorrect: false},
					{Text: fmt.Sprintf("x = %d", i/2), IsCorrect: false},
					{Text: fmt.Sprintf("x = -%d", i), IsCorrect: false},
				},
			}
		case 5:
			q = QuestionData{
				Text: fmt.Sprintf("What is the modulus of the complex number z = %d + 4i?", 3),
				Topic: topic,
				Explanation: "|z| = √(3² + 4²) = √(9 + 16) = 5.",
				Options: []Option{
					{Text: "5", IsCorrect: true},
					{Text: "7", IsCorrect: false},
					{Text: "25", IsCorrect: false},
					{Text: "√7", IsCorrect: false},
				},
			}
		case 6:
			q = QuestionData{
				Text: fmt.Sprintf("If vector A = (%d, 0, 0) and B = (0, %d, 0), what is the cross product A × B?", i%5+1, i%5+1),
				Topic: topic,
				Explanation: fmt.Sprintf("A × B = (0, 0, %d).", (i%5+1)*(i%5+1)),
				Options: []Option{
					{Text: fmt.Sprintf("(0, 0, %d)", (i%5+1)*(i%5+1)), IsCorrect: true},
					{Text: "(0, 0, 0)", IsCorrect: false},
					{Text: fmt.Sprintf("(%d, %d, 0)", i%5+1, i%5+1), IsCorrect: false},
					{Text: fmt.Sprintf("(0, %d, 0)", (i%5+1)*(i%5+1)), IsCorrect: false},
				},
			}
		case 7:
			q = QuestionData{
				Text: fmt.Sprintf("What is the variance of a fair standard %d-sided die (from 1 to 6)?", 6),
				Topic: topic,
				Explanation: "For uniform discrete 1..6, variance = (n² - 1)/12 = (36 - 1)/12 = 35/12 ≈ 2.92.",
				Options: []Option{
					{Text: "35 / 12", IsCorrect: true},
					{Text: "7 / 2", IsCorrect: false},
					{Text: "5 / 12", IsCorrect: false},
					{Text: "15 / 4", IsCorrect: false},
				},
			}
		default:
			q = QuestionData{
				Text: fmt.Sprintf("Evaluate the limit as x -> ∞ of (%dx² + 5) / (%dx² - 3x).", i, i*2),
				Topic: topic,
				Explanation: fmt.Sprintf("Ratio of leading coefficients = %d / %d = 1/2.", i, i*2),
				Options: []Option{
					{Text: "1/2", IsCorrect: true},
					{Text: "1", IsCorrect: false},
					{Text: "0", IsCorrect: false},
					{Text: "∞", IsCorrect: false},
				},
			}
		}
		questions = append(questions, q)
	}

	return ExamData{
		Title:           "Grade 12 National Mathematics Mock Examination",
		Subject:         "Mathematics",
		Grade:           12,
		DurationMinutes: 120,
		PassMarks:       30,
		Instructions:    "This examination contains 60 multiple-choice questions. Select the most correct answer for each problem. Calculators are not allowed.",
		Questions:       questions,
	}
}
