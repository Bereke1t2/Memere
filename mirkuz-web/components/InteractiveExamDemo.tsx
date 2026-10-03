"use client";

import React, { useState } from "react";
import { 
  CheckCircle2, 
  XCircle, 
  Sparkles, 
  Clock, 
  HelpCircle, 
  ArrowRight, 
  RefreshCw,
  Lightbulb,
  Check
} from "lucide-react";

interface Question {
  id: number;
  subject: string;
  year: string;
  prompt: string;
  formula?: string;
  options: { label: string; text: string }[];
  correctAnswer: string;
  explanation: string;
  conceptTag: string;
}

const sampleQuestions: Question[] = [
  {
    id: 1,
    subject: "Mathematics (Natural)",
    year: "2016 E.C. National Entrance Exam",
    prompt: "What is the derivative of the function f(x) with respect to x?",
    formula: "f(x) = \\ln(x^2 + 4x + 5)",
    options: [
      { label: "A", text: "f'(x) = (2x + 4) / (x^2 + 4x + 5)" },
      { label: "B", text: "f'(x) = 1 / (x^2 + 4x + 5)" },
      { label: "C", text: "f'(x) = (2x + 4) * \\ln(x^2 + 4x + 5)" },
      { label: "D", text: "f'(x) = (x + 2) / (x^2 + 4x + 5)" },
    ],
    correctAnswer: "A",
    explanation:
      "By the chain rule for logarithmic differentiation: d/dx[ln(u)] = u'/u. Let u = x^2 + 4x + 5. Then u' = 2x + 4. Therefore, f'(x) = (2x + 4)/(x^2 + 4x + 5). Correct option is A.",
    conceptTag: "Calculus — Chain Rule of Logarithmic Functions",
  },
  {
    id: 2,
    subject: "Physics (Natural)",
    year: "2015 E.C. National Entrance Exam",
    prompt:
      "A projectile is launched from ground level with an initial velocity of 50 m/s at an angle of 30° above the horizontal. (Take g = 10 m/s²). What is the maximum height reached by the projectile?",
    formula: "H_{max} = \\frac{u^2 \\sin^2(\\theta)}{2g}",
    options: [
      { label: "A", text: "125.0 m" },
      { label: "B", text: "31.25 m" },
      { label: "C", text: "62.5 m" },
      { label: "D", text: "15.62 m" },
    ],
    correctAnswer: "B",
    explanation:
      "Maximum height formula: H = (u * sin(θ))² / (2g). Here u = 50 m/s, θ = 30° => sin(30°) = 0.5. Vertical velocity u_y = 50 * 0.5 = 25 m/s. H = (25)² / (2 * 10) = 625 / 20 = 31.25 m. Correct option is B.",
    conceptTag: "Mechanics — Two-Dimensional Projectile Motion",
  },
  {
    id: 3,
    subject: "Aptitude & Logic",
    year: "2016 E.C. National Entrance Exam",
    prompt:
      "If all Zors are Blips, and some Blips are Quarks, but no Quarks are Flips, which of the following statements must be conclusively true?",
    options: [
      { label: "A", text: "All Zors are Quarks" },
      { label: "B", text: "Some Blips are definitely not Flips" },
      { label: "C", text: "No Zors are Flips" },
      { label: "D", text: "All Flips are Zors" },
    ],
    correctAnswer: "B",
    explanation:
      "Since some Blips are Quarks, and no Quarks are Flips, those specific Blips that are Quarks cannot possibly be Flips. Therefore, 'Some Blips are definitely not Flips' is logically valid and guaranteed true.",
    conceptTag: "Analytical Reasoning — Categorical Syllogisms",
  },
];

export function InteractiveExamDemo() {
  const [selectedQuestionIdx, setSelectedQuestionIdx] = useState(0);
  const [selectedAnswer, setSelectedAnswer] = useState<string | null>(null);
  const [showExplanation, setShowExplanation] = useState(false);
  const [timerSeconds, setTimerSeconds] = useState(120);

  const currentQ = sampleQuestions[selectedQuestionIdx];

  const handleSelectAnswer = (label: string) => {
    setSelectedAnswer(label);
    setShowExplanation(true);
  };

  const handleReset = (idx: number) => {
    setSelectedQuestionIdx(idx);
    setSelectedAnswer(null);
    setShowExplanation(false);
  };

  return (
    <section id="exam-demo" className="bg-[#F5F7FA] py-20 border-b border-[#E8ECF2]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Section Header */}
        <div className="text-center max-w-3xl mx-auto space-y-2 mb-12">
          <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-[#E8F5E9] text-[#4CAF4F] text-xs font-semibold">
            <Sparkles className="w-3.5 h-3.5" />
            <span>Interactive Question Simulator</span>
          </div>
          <h2 className="text-3xl sm:text-4xl font-bold text-[#4D4D4D] tracking-tight">
            Try a Real National Entrance Exam Question
          </h2>
          <p className="text-sm sm:text-base text-[#717171]">
            Experience instant step-by-step video solutions, formulas, and timing just like the official Ministry of Education exam.
          </p>
        </div>

        {/* Interactive Box */}
        <div className="max-w-4xl mx-auto bg-white rounded-xl border border-[#E8ECF2] shadow-[0_4px_12px_rgba(171,190,209,0.25)] overflow-hidden">
          
          {/* Top Bar / Question Selector */}
          <div className="bg-[#F5F7FA] border-b border-[#E8ECF2] p-4 sm:px-6 flex flex-wrap items-center justify-between gap-4">
            <div className="flex items-center gap-2 overflow-x-auto">
              {sampleQuestions.map((q, idx) => (
                <button
                  key={q.id}
                  onClick={() => handleReset(idx)}
                  className={`px-3 py-1.5 rounded-md text-xs font-semibold transition-all shrink-0 ${
                    selectedQuestionIdx === idx
                      ? "bg-[#4CAF4F] text-white shadow-xs"
                      : "bg-white text-[#717171] border border-[#E8ECF2] hover:text-[#4D4D4D]"
                  }`}
                >
                  Q{q.id}: {q.subject.split(" ")[0]}
                </button>
              ))}
            </div>

            <div className="flex items-center gap-3 text-xs text-[#717171]">
              <div className="flex items-center gap-1.5 px-2.5 py-1 rounded bg-white border border-[#E8ECF2]">
                <Clock className="w-3.5 h-3.5 text-[#4CAF4F]" />
                <span className="font-mono font-bold text-[#4D4D4D]">01:45</span>
              </div>
              <span className="text-[11px] font-medium text-[#4CAF4F] bg-[#E8F5E9] px-2 py-0.5 rounded">
                {currentQ.year}
              </span>
            </div>
          </div>

          {/* Question Content */}
          <div className="p-6 sm:p-8 space-y-6">
            
            {/* Subject & Concept tag */}
            <div className="flex items-center justify-between text-xs">
              <span className="font-bold text-[#4D4D4D]">{currentQ.subject}</span>
              <span className="text-[#89939E] italic">{currentQ.conceptTag}</span>
            </div>

            {/* Prompt */}
            <div className="space-y-3">
              <p className="text-base sm:text-lg font-medium text-[#263238] leading-relaxed">
                {currentQ.prompt}
              </p>
              {currentQ.formula && (
                <div className="p-3 rounded-lg bg-[#F5F7FA] font-mono text-sm text-[#4CAF4F] border border-[#E8ECF2] inline-block font-semibold">
                  {currentQ.formula}
                </div>
              )}
            </div>

            {/* Options */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2">
              {currentQ.options.map((opt) => {
                const isSelected = selectedAnswer === opt.label;
                const isCorrect = opt.label === currentQ.correctAnswer;
                let btnStyle = "border-[#E8ECF2] bg-white hover:border-[#4CAF4F] text-[#4D4D4D]";

                if (showExplanation) {
                  if (isCorrect) {
                    btnStyle = "border-[#4CAF4F] bg-[#E8F5E9] text-[#2E7D32] font-semibold";
                  } else if (isSelected && !isCorrect) {
                    btnStyle = "border-red-400 bg-red-50 text-red-700";
                  } else {
                    btnStyle = "border-[#E8ECF2] bg-white text-[#89939E] opacity-60";
                  }
                }

                return (
                  <button
                    key={opt.label}
                    onClick={() => handleSelectAnswer(opt.label)}
                    disabled={showExplanation}
                    className={`flex items-start gap-3 p-4 rounded-lg border text-left text-sm transition-all ${btnStyle}`}
                  >
                    <span className="w-6 h-6 rounded-full bg-[#F5F7FA] border border-[#E8ECF2] flex items-center justify-center font-bold text-xs shrink-0 mt-0.5">
                      {opt.label}
                    </span>
                    <span className="flex-1 leading-snug">{opt.text}</span>
                    {showExplanation && isCorrect && (
                      <CheckCircle2 className="w-5 h-5 text-[#4CAF4F] shrink-0" />
                    )}
                    {showExplanation && isSelected && !isCorrect && (
                      <XCircle className="w-5 h-5 text-red-500 shrink-0" />
                    )}
                  </button>
                );
              })}
            </div>

            {/* Explanation box */}
            {showExplanation && (
              <div className="p-5 rounded-lg bg-[#E8F5E9] border border-[#C8E6C9] space-y-2 animate-in fade-in duration-300">
                <div className="flex items-center gap-2 text-[#2E7D32] font-bold text-sm">
                  <Lightbulb className="w-4 h-4" />
                  <span>Official Step-by-Step Solution:</span>
                </div>
                <p className="text-xs sm:text-sm text-[#1B5E20] leading-relaxed">
                  {currentQ.explanation}
                </p>
              </div>
            )}

            {/* Bottom action reset/next */}
            <div className="flex items-center justify-between pt-2">
              <button
                onClick={() => handleReset(selectedQuestionIdx)}
                className="inline-flex items-center gap-1.5 text-xs text-[#717171] hover:text-[#4D4D4D] transition-colors"
              >
                <RefreshCw className="w-3.5 h-3.5" />
                <span>Reset Question</span>
              </button>

              <button
                onClick={() => handleReset((selectedQuestionIdx + 1) % sampleQuestions.length)}
                className="inline-flex items-center gap-2 px-5 py-2 rounded-md bg-[#4CAF4F] hover:bg-[#388E3C] text-white text-xs font-semibold shadow-xs transition-all"
              >
                <span>Next Question</span>
                <ArrowRight className="w-3.5 h-3.5" />
              </button>
            </div>

          </div>

        </div>

      </div>
    </section>
  );
}
