package main

func buildHistoryCivicsCourse() CourseData {
	return CourseData{
		Title:            "Grade 12 History & Civic Education: National University Entrance Prep",
		Slug:             "grade-12-history-civics-euee-prep",
		Subject:          "Civics",
		Grade:            12,
		Description:      "Comprehensive preparation for Ethiopian National Examination in Modern Ethiopian & World History, Constitutional Development, Human Rights, Democratic Governance, and Civic Responsibilities. Includes in-depth historical analyses, past entrance exam questions, unit quizzes, and mock exams.",
		ShortDescription: "Complete Grade 12 History & Civics EUEE prep covering Ethiopian modern history, Adwa, constitution, democracy, quizzes, and exams.",
		ThumbnailUrl:     "https://images.unsplash.com/photo-1461360370896-922624d12aa1?w=800&auto=format&fit=crop&q=80",
		Price:            0,
		IsFree:           true,
		Level:            "intermediate",
		Sections: []SectionData{
			{
				Title:       "Unit 1: Modern Ethiopian State Formation & The Battle of Adwa",
				Description: "Reunification of the Empire under Tewodros II, Yohannes IV, Menelik II, Treaties of Wuchale, and the victory of Adwa (1896).",
				Lessons: []LessonData{
					{
						Title:           "Lesson 1.1: The Era of Princes (Zemene Mesafint) to Centralization",
						Type:            "video",
						DurationSeconds: 1200,
						IsFreePreview:   true,
						Content: `# Modern Ethiopian State Formation
The modern Ethiopian state emerged from the fragmentation of the **Zemene Mesafint** (1769–1855).

### Key Historical Figures:
1. **Emperor Tewodros II (1855–1868)**:
   - Initiated modern military reforms, centralized tax collection, and challenged regional lords.
   - Ended with the British Napier Expedition at Meqdela (1868).
2. **Emperor Yohannes IV (1872–1889)**:
   - Defended Ethiopian sovereignty against Egyptian expansionism (Gundet 1875, Gura 1876), Italian incursions (Dogali 1887), and Mahdist forces (Metemma 1889).
3. **Emperor Menelik II (1889–1913)**:
   - Unified southern territories, modernized infrastructure (railway, telephone, postal system, modern schools).`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 1.2: The Battle of Adwa (1896) & Global Impact",
						Type:            "note",
						DurationSeconds: 900,
						IsFreePreview:   true,
						Content: `# The Battle of Adwa (March 1, 1896)
- **Treaty of Wuchale (1889)**: Article XVII dispute (Italian version claimed an Italian protectorate over Ethiopia, while Amharic version preserved independence).
- Ethiopian forces under Emperor Menelik II and Empress Taytu decisively defeated the Italian colonial army commanded by General Baratieri.
- **Treaty of Addis Ababa (October 1896)**: Unconditional annulment of the Treaty of Wuchale and absolute recognition of Ethiopian independence by Italy.
- **Global Impact**: Inspired anti-colonial movements across Africa, the African diaspora, and the Pan-African movement.`,
						PdfUrl: "sample.pdf",
					},
				},
			},
			{
				Title:       "Unit 2: World Wars, League of Nations & Italian Fascist Occupation",
				Description: "Origins of WWI & WWII, the Walwal incident, Italian fascist invasion (1935–1941), Patriot resistance, and Liberation.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 2.1: The 1935–1941 Fascist Invasion & Patriotic Resistance",
						Type:            "video",
						DurationSeconds: 1350,
						IsFreePreview:   false,
						Content: `# Italian Fascist Occupation (1935–1941)
- **Walwal Incident (Dec 1934)**: Fabricated border clash used by Mussolini as pretext for invasion.
- Invasion began in October 1935 using illegal chemical mustard gas against Ethiopian armies.
- **The Black Lions (Tikur Anbessa)**: Educated military officers leading guerrilla warfare.
- **Patriots (Arbegnoch)**: Leaders like Ras Desta Damtew, Dejazmach Balcha Safo, and Belay Zeleke.
- **Yekatit 12 Massacre (1937)**: Following the assassination attempt on Graziani by Abrha Deboch and Moges Asgedom.
- **Liberation (May 5, 1941)**: Victory with the Gideon Force and entry of Emperor Haile Selassie into Addis Ababa.`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 2.2: The League of Nations & Post-WWII International Order",
						Type:            "note",
						DurationSeconds: 850,
						IsFreePreview:   false,
						Content: `# International Organizations & Collective Security
- Failure of the **League of Nations** to enforce sanctions against fascist aggression exposed the weakness of collective security.
- Creation of the **United Nations (UN)** in 1945 with Ethiopia as a founding member.
- Formation of the **Organization of African Unity (OAU)** in Addis Ababa (May 25, 1963).`,
						PdfUrl: "sample.pdf",
					},
				},
			},
			{
				Title:       "Unit 3: Constitutional Development in Ethiopia",
				Description: "The 1931 Constitution, 1955 Revised Constitution, 1987 PDRE Constitution, and the 1995 FDRE Constitution.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 3.1: Historical Evolution of Ethiopian Constitutions",
						Type:            "video",
						DurationSeconds: 1250,
						IsFreePreview:   false,
						Content: `# Evolution of Written Constitutions in Ethiopia
1. **1931 Constitution**: First written constitution in Ethiopian history; concentrated supreme executive, legislative, and judicial power in the Emperor.
2. **1955 Revised Constitution**: Modernized bicameral parliament, introduced bill of rights with severe limitations, reinforced imperial succession.
3. **1987 PDRE Constitution**: Socialist single-party state under the Workers' Party of Ethiopia (WPE).
4. **1995 FDRE Constitution**: Established the Federal Democratic Republic with ethnic federalism, bicameral parliament (HPR & HoF), and self-determination rights (Article 39).`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 3.2: Separation of Powers & The Rule of Law",
						Type:            "note",
						DurationSeconds: 850,
						IsFreePreview:   false,
						Content: `# Principles of Democratic Governance
- **Rule of Law**: All individuals and state institutions are accountable to publicly promulgated laws.
- **Separation of Powers**:
  - **Legislature**: Enacts laws (House of Peoples' Representatives - HPR).
  - **Executive**: Implements policy (Prime Minister and Council of Ministers).
  - **Judiciary**: Interprets laws (Federal Supreme Court).
- **House of Federation (HoF)**: Interprets the constitution and handles regional subsidies.`,
						PdfUrl: "sample.pdf",
					},
				},
			},
			{
				Title:       "Unit 4: Human Rights, Democracy & Citizen Participation",
				Description: "Generations of human rights, democratic values, civic participation, peaceful dispute resolution, and national development.",
				Lessons: []LessonData{
					{
						Title:           "Lesson 4.1: Three Generations of Human Rights",
						Type:            "video",
						DurationSeconds: 1150,
						IsFreePreview:   false,
						Content: `# Classification of Human Rights
1. **First-Generation (Civil & Political Rights)**:
   - Right to life, liberty, fair trial, freedom of speech, assembly, and voting (Negative rights).
2. **Second-Generation (Economic, Social & Cultural Rights)**:
   - Right to education, healthcare, work, and adequate standard of living (Positive rights).
3. **Third-Generation (Solidarity / Collective Rights)**:
   - Right to peace, clean environment, development, and self-determination.`,
						PdfUrl: "sample.pdf",
					},
					{
						Title:           "Lesson 4.2: Civic Ethics, Patriotism & Conflict Resolution",
						Type:            "note",
						DurationSeconds: 800,
						IsFreePreview:   false,
						Content: `# Civic Engagement & Conflict Management
- **Democratic Patriotism**: Active participation in societal progress, paying taxes, obeying the law, defending national sovereignty, and promoting social harmony.
- **Alternative Dispute Resolution (ADR)**: Negotiation, Mediation, and Arbitration vs. Litigation.
- Traditional Ethiopian conflict resolution institutions (e.g. *Shimgellina*, *Gadaa/Jaarsummaa*, *Bayto*).`,
						PdfUrl: "sample.pdf",
					},
				},
			},
		},
		Quizzes: []QuizData{
			{
				Title:            "History Quiz 1: 19th Century State Formation & Adwa",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "Which article of the 1889 Treaty of Wuchale contained the controversial mistranslation that triggered the Battle of Adwa?",
						Explanation: "Article XVII differed between the Amharic and Italian texts, where the Italian version falsely claimed an Italian protectorate over Ethiopia.",
						Options: []Option{
							{Text: "Article XVII (17)", IsCorrect: true},
							{Text: "Article III (3)", IsCorrect: false},
							{Text: "Article X (10)", IsCorrect: false},
							{Text: "Article XII (12)", IsCorrect: false},
						},
					},
					{
						Text:        "In which year did the historic Battle of Adwa take place?",
						Explanation: "The Battle of Adwa took place on March 1, 1896 (Yekatit 23, 1888 E.C.).",
						Options: []Option{
							{Text: "1896", IsCorrect: true},
							{Text: "1889", IsCorrect: false},
							{Text: "1935", IsCorrect: false},
							{Text: "1868", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "History Quiz 2: Fascist Occupation & Patriotic Resistance",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "Who were the two young patriots responsible for the historic assassination attempt on Marshal Graziani in Addis Ababa in February 1937?",
						Explanation: "Abrha Deboch and Moges Asgedom threw grenades at Marshal Graziani at the Genete Leul Palace on Yekatit 12, 1937.",
						Options: []Option{
							{Text: "Abrha Deboch and Moges Asgedom", IsCorrect: true},
							{Text: "Ras Desta Damtew and Dejazmach Balcha", IsCorrect: false},
							{Text: "Belay Zeleke and Jagama Kello", IsCorrect: false},
							{Text: "Geresu Duki and Amoraw Wubneh", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Civics Quiz 3: Ethiopian Constitutional Development",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "What was the main purpose of the first written Ethiopian Constitution promulgated in 1931?",
						Explanation: "The 1931 constitution concentrated absolute state power in the hands of Emperor Haile Selassie and created a loyal centralized administrative hierarchy.",
						Options: []Option{
							{Text: "Centralize power in the Emperor and formalize imperial succession", IsCorrect: true},
							{Text: "Establish multi-party democratic governance", IsCorrect: false},
							{Text: "Institute a socialist republic", IsCorrect: false},
							{Text: "Grant regional autonomy and self-determination", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Civics Quiz 4: Organs of Democratic Governance",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "Under the 1995 FDRE Constitution, which government institution is constitutionally vested with the power to interpret the Constitution?",
						Explanation: "According to the FDRE Constitution, the House of Federation (HoF), advised by the Council of Constitutional Inquiry, has the mandate to interpret the constitution.",
						Options: []Option{
							{Text: "House of Federation (HoF)", IsCorrect: true},
							{Text: "Federal Supreme Court", IsCorrect: false},
							{Text: "House of Peoples' Representatives (HPR)", IsCorrect: false},
							{Text: "Prime Minister's Office", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:            "Civics Quiz 5: Human Rights & Citizen Ethics",
				TimeLimitSeconds: 600,
				PassPercentage:   70.0,
				Questions: []QuestionData{
					{
						Text:        "Civil and political liberties such as freedom of expression and the right to a fair trial belong to which generation of human rights?",
						Explanation: "Civil and political rights are categorized as First-Generation human rights (liberty-oriented).",
						Options: []Option{
							{Text: "First-Generation Rights", IsCorrect: true},
							{Text: "Second-Generation Rights", IsCorrect: false},
							{Text: "Third-Generation Rights", IsCorrect: false},
							{Text: "Fourth-Generation Rights", IsCorrect: false},
						},
					},
				},
			},
		},
		Exams: []ExamData{
			{
				Title:           "History Unit 1 & 2 Modern Ethiopian History Exam",
				Subject:         "Civics",
				Grade:           12,
				DurationMinutes: 45,
				PassMarks:       12,
				Instructions:    "Covers 19th Century State Formation, Treaty of Wuchale, Battle of Adwa, and 1935-1941 Patriot Resistance. 45 minutes.",
				Questions: []QuestionData{
					{
						Text:        "Which European military expedition led to the tragic death of Emperor Tewodros II at Meqdela in 1868?",
						Explanation: "The British expeditionary force commanded by Sir Robert Napier attacked Meqdela in April 1868.",
						Options: []Option{
							{Text: "The British expedition led by Robert Napier", IsCorrect: true},
							{Text: "The Italian colonial army led by Baratieri", IsCorrect: false},
							{Text: "The Egyptian expedition led by Werner Munzinger", IsCorrect: false},
							{Text: "The French army led by Marchand", IsCorrect: false},
						},
					},
					{
						Text:        "Which international body was founded in Addis Ababa in May 1963 to promote solidarity among newly independent African states?",
						Explanation: "The Organization of African Unity (OAU) was founded on May 25, 1963 in Addis Ababa.",
						Options: []Option{
							{Text: "Organization of African Unity (OAU)", IsCorrect: true},
							{Text: "Economic Community of West African States (ECOWAS)", IsCorrect: false},
							{Text: "Non-Aligned Movement (NAM)", IsCorrect: false},
							{Text: "League of Arab States", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:           "Civics Midterm Examination (Constitutions & Democratic Institutions)",
				Subject:         "Civics",
				Grade:           12,
				DurationMinutes: 60,
				PassMarks:       15,
				Instructions:    "Testing Constitutional History, Separation of Powers, Rule of Law, and Federalism.",
				Questions: []QuestionData{
					{
						Text:        "What distinguishes First-Generation human rights from Second-Generation human rights?",
						Explanation: "First-generation rights focus on civil and political liberties (non-interference by the state), while second-generation rights focus on socioeconomic entitlements (education, healthcare, housing).",
						Options: []Option{
							{Text: "First-generation focuses on civil liberties; second-generation on socioeconomic entitlements", IsCorrect: true},
							{Text: "First-generation focuses on environmental rights; second-generation on political voting", IsCorrect: false},
							{Text: "First-generation applies only to corporations; second-generation to individuals", IsCorrect: false},
							{Text: "First-generation is non-binding; second-generation is strictly enforced", IsCorrect: false},
						},
					},
				},
			},
			{
				Title:           "Ethiopian National EUEE History & Civics Final Mock Exam",
				Subject:         "Civics",
				Grade:           12,
				DurationMinutes: 90,
				PassMarks:       25,
				Instructions:    "Comprehensive National Entrance Examination simulation for Grade 12 History & Civic Education.",
				Questions: []QuestionData{
					{
						Text:        "Which peaceful dispute resolution method involves a neutral third party facilitating dialogue without imposing a binding decision?",
						Explanation: "Mediation involves a neutral mediator assisting disputing parties to reach a mutually agreeable voluntary settlement.",
						Options: []Option{
							{Text: "Mediation", IsCorrect: true},
							{Text: "Arbitration", IsCorrect: false},
							{Text: "Litigation", IsCorrect: false},
							{Text: "Adjudication", IsCorrect: false},
						},
					},
				},
			},
		},
	}
}
