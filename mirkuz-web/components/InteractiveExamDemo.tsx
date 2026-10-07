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
    <section id="exam-demo" className="bg-white py-20 md:py-28 border-b border-[#E2E8F0]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Section Header */}
        <div className="text-center max-w-3xl mx-auto space-y-3 mb-14">
          <div className="inline-flex items-center gap-2 px-3.5 py-1 rounded-full bg-[#E8F8F5] text-[#00B894] text-xs font-semibold">
            <Sparkles className="w-3.5 h-3.5" />
            <span>Interactive Question Simulator</span>
          </div>
          <h2 className="text-3xl sm:text-4xl font-extrabold text-[#1E2B58] tracking-tight">
            Try a Real National <span className="title-underline-center text-[#1E2B58]">Entrance Exam Question</span>
          </h2>
          <p className="text-sm sm:text-base text-[#64748B] pt-1">
            Experience instant step-by-step video solutions, formulas, and timing just like the official Ministry of Education exam.
          </p>
        </div>

        {/* Interactive Box */}
        <div className="max-w-4xl mx-auto bg-white rounded-3xl border border-[#E2E8F0] shadow-dribbble-lg overflow-hidden">
          
          {/* Top Bar / Question Selector */}
          <div className="bg-[#F8FAFC] border-b border-[#E2E8F0] p-4 sm:px-6 flex flex-wrap items-center justify-between gap-4">
            <div className="flex items-center gap-2 overflow-x-auto">
              {sampleQuestions.map((q, idx) => (
                <button
                  key={q.id}
                  onClick={() => handleReset(idx)}
                  className={`px-3.5 py-1.5 rounded-full text-xs font-bold transition-all shrink-0 ${
                    selectedQuestionIdx === idx
                      ? "bg-[#1E2B58] text-white shadow-xs"
                      : "bg-white text-[#64748B] border border-[#E2E8F0] hover:text-[#1E2B58]"
                  }`}
                >
                  Q{q.id}: {q.subject.split(" ")[0]}
                </button>
              ))}
            </div>

            <div className="flex items-center gap-3 text-xs text-[#64748B]">
              <div className="flex items-center gap-1.5 px-3 py-1 rounded-full bg-[#E8F8F5] border border-[#00B894]/20">
                <Clock className="w-3.5 h-3.5 text-[#00B894]" />
                <span className="font-mono font-bold text-[#1E2B58]">01:45</span>
              </div>
              <span className="text-[11px] font-semibold text-[#1E2B58] bg-[#1E2B58]/10 px-2.5 py-1 rounded-full">
                {currentQ.year}
              </span>
            </div>
          </div>

          {/* Question Content */}
          <div className="p-6 sm:p-8 space-y-6">
            
            {/* Subject & Concept tag */}
            <div className="flex items-center justify-between text-xs">
              <span className="font-bold text-[#1E2B58]">{currentQ.subject}</span>
              <span className="text-[#94A3B8] italic">{currentQ.conceptTag}</span>
            </div>

            {/* Prompt */}
            <div className="space-y-3">
              <p className="text-base sm:text-lg font-semibold text-[#1E2B58] leading-relaxed">
                {currentQ.prompt}
              </p>
              {currentQ.formula && (
                <div className="p-3 rounded-xl bg-[#F8FAFC] font-mono text-sm text-[#00B894] border border-[#E2E8F0] inline-block font-bold">
                  {currentQ.formula}
                </div>
              )}
            </div>

            {/* Options */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2">
              {currentQ.options.map((opt) => {
                const isSelected = selectedAnswer === opt.label;
                const isCorrect = opt.label === currentQ.correctAnswer;
                let btnStyle = "border-[#E2E8F0] bg-white hover:border-[#00B894] text-[#1E2B58] hover:bg-[#F8FAFC]";

                if (showExplanation) {
                  if (isCorrect) {
                    btnStyle = "border-[#00B894] bg-[#E8F8F5] text-[#065F46] font-semibold";
                  } else if (isSelected && !isCorrect) {
                    btnStyle = "border-rose-400 bg-rose-50 text-rose-700";
                  } else {
                    btnStyle = "border-[#E2E8F0] bg-white text-[#94A3B8] opacity-60";
                  }
                }

                return (
                  <button
                    key={opt.label}
                    onClick={() => handleSelectAnswer(opt.label)}
                    disabled={showExplanation}
                    className={`flex items-start gap-3 p-4 rounded-xl border text-left text-sm transition-all cursor-pointer ${btnStyle}`}
                  >
                    <span className="w-6 h-6 rounded-full bg-[#F1F5F9] border border-[#E2E8F0] flex items-center justify-center font-bold text-xs shrink-0 mt-0.5 text-[#1E2B58]">
                      {opt.label}
                    </span>
                    <span className="flex-1 leading-snug">{opt.text}</span>
                    {showExplanation && isCorrect && (
                      <CheckCircle2 className="w-5 h-5 text-[#00B894] shrink-0" />
                    )}
                    {showExplanation && isSelected && !isCorrect && (
                      <XCircle className="w-5 h-5 text-rose-500 shrink-0" />
                    )}
                  </button>
                );
              })}
            </div>

            {/* Explanation box */}
            {showExplanation && (
              <div className="p-5 rounded-2xl bg-[#E8F8F5] border border-[#00B894]/30 space-y-2 animate-in fade-in duration-300">
                <div className="flex items-center gap-2 text-[#065F46] font-bold text-sm">
                  <Lightbulb className="w-4 h-4 text-[#00B894]" />
                  <span>Official Step-by-Step Solution:</span>
                </div>
                <p className="text-xs sm:text-sm text-[#065F46] leading-relaxed">
                  {currentQ.explanation}
                </p>
              </div>
            )}

            {/* Bottom action reset/next */}
            <div className="flex items-center justify-between pt-2">
              <button
                onClick={() => handleReset(selectedQuestionIdx)}
                className="inline-flex items-center gap-1.5 text-xs text-[#64748B] hover:text-[#1E2B58] transition-colors cursor-pointer"
              >
                <RefreshCw className="w-3.5 h-3.5" />
                <span>Reset Question</span>
              </button>

              <button
                onClick={() => handleReset((selectedQuestionIdx + 1) % sampleQuestions.length)}
                className="inline-flex items-center gap-2 px-6 py-2.5 rounded-full bg-[#00B894] hover:bg-[#009874] text-white text-xs font-bold shadow-xs transition-all cursor-pointer"
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
