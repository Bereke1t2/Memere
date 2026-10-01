"use client";

import React, { useState } from "react";
import { 
  CheckCircle2, 
  XCircle, 
  Sparkles, 
  Clock, 
  HelpCircle, 
  ArrowRight, 
  Download, 
  Send,
  Lightbulb,
  Check,
  RefreshCw
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
      { label: "A", text: "f'(x) = \\frac{2x + 4}{x^2 + 4x + 5}" },
      { label: "B", text: "f'(x) = \\frac{1}{x^2 + 4x + 5}" },
      { label: "C", text: "f'(x) = (2x + 4) \\ln(x^2 + 4x + 5)" },
      { label: "D", text: "f'(x) = \\frac{x + 2}{x^2 + 4x + 5}" },
    ],
    correctAnswer: "A",
    explanation:
      "By applying the Chain Rule for natural logarithms: \\frac{d}{dx}[\\ln(u)] = \\frac{u'}{u}. Here, let u = x^2 + 4x + 5, so u' = 2x + 4. Therefore, f'(x) = \\frac{2x + 4}{x^2 + 4x + 5}.",
    conceptTag: "Unit 3: Differential Calculus",
  },
  {
    id: 2,
    subject: "Physics",
    year: "2015 E.C. National Entrance Exam",
    prompt: "An object is thrown vertically upward with an initial velocity of 20 m/s. Neglecting air resistance (g = 10 m/s²), what is the maximum height reached by the object?",
    options: [
      { label: "A", text: "10 meters" },
      { label: "B", text: "20 meters" },
      { label: "C", text: "40 meters" },
      { label: "D", text: "50 meters" },
    ],
    correctAnswer: "B",
    explanation:
      "Using the kinematic equation v^2 = u^2 - 2gh. At maximum height, final velocity v = 0 m/s. Thus, 0 = (20)^2 - 2(10)h  =>  20h = 400  =>  h = 20 meters.",
    conceptTag: "Unit 2: Kinematics & Vertical Motion",
  },
];

export function InteractiveExamDemo() {
  const [selectedQuestionIndex, setSelectedQuestionIndex] = useState(0);
  const [selectedOption, setSelectedOption] = useState<string | null>(null);
  const [showExplanation, setShowExplanation] = useState(false);

  const currentQ = sampleQuestions[selectedQuestionIndex];

  const handleSelectOption = (label: string) => {
    setSelectedOption(label);
    setShowExplanation(true);
  };

  const handleNextQuestion = () => {
    setSelectedOption(null);
    setShowExplanation(false);
    setSelectedQuestionIndex((prev) => (prev + 1) % sampleQuestions.length);
  };

  const isCorrect = selectedOption === currentQ.correctAnswer;

  return (
    <section id="exam-demo" className="py-24 bg-[#050505] relative overflow-hidden">
      {/* Background glow */}
      <div className="absolute top-1/3 left-1/2 -translate-x-1/2 w-[600px] h-[400px] bg-[#10b981]/10 rounded-full blur-[140px] pointer-events-none" />

      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        
        {/* Section Title */}
        <div className="text-center mb-12 space-y-3">
          <div className="inline-flex items-center gap-2 px-3.5 py-1 rounded-full bg-[#10b981]/15 text-[#10b981] border border-[#10b981]/30 text-xs font-bold">
            <Sparkles className="w-3.5 h-3.5" />
            <span>Interactive Live Micro-Demo</span>
          </div>
          <h2 className="text-3xl sm:text-4xl font-extrabold text-white tracking-tight">
            Try a Real National Entrance Exam Question
          </h2>
          <p className="text-sm sm:text-base text-[#94a3b8] max-w-xl mx-auto">
            Experience how Mirkuz turns tricky exam problems into clear, memorable concepts with instant explanations.
          </p>

          {/* Question Switcher Tabs */}
          <div className="flex justify-center gap-2 pt-2">
            {sampleQuestions.map((q, idx) => (
              <button
                key={q.id}
                onClick={() => {
                  setSelectedQuestionIndex(idx);
                  setSelectedOption(null);
                  setShowExplanation(false);
                }}
                className={`px-3.5 py-1.5 rounded-lg text-xs font-semibold transition-all ${
                  selectedQuestionIndex === idx
                    ? "bg-[#18181e] text-white border border-[#10b981]"
                    : "bg-[#111115] text-[#94a3b8] border border-[#1a1a20] hover:text-white"
                }`}
              >
                Sample {idx + 1}: {q.subject}
              </button>
            ))}
          </div>
        </div>

        {/* Exam Card Container */}
        <div className="rounded-3xl bg-[#111115] border border-[#2c2c35] p-6 sm:p-8 shadow-2xl shadow-black/80 space-y-6">
          
          {/* Header row: Subject, Year & Mock Timer */}
          <div className="flex flex-wrap items-center justify-between gap-3 pb-4 border-b border-[#1a1a20] text-xs">
            <div className="flex items-center gap-2">
              <span className="px-2.5 py-1 rounded-md bg-[#18181e] text-[#10b981] font-bold border border-[#10b981]/25">
                {currentQ.subject}
              </span>
              <span className="text-[#94a3b8] font-medium">{currentQ.year}</span>
            </div>
            <div className="flex items-center gap-1.5 text-[#f59e0b] font-mono font-bold bg-[#f59e0b]/10 px-2.5 py-1 rounded-md border border-[#f59e0b]/20">
              <Clock className="w-3.5 h-3.5" />
              <span>Exam Clock: 01:28</span>
            </div>
          </div>

          {/* Question Prompt */}
          <div className="space-y-3">
            <div className="text-[11px] font-bold uppercase tracking-wider text-[#38bdf8]">
              {currentQ.conceptTag}
            </div>
            <p className="text-base sm:text-lg font-medium text-white leading-relaxed">
              {currentQ.prompt}
            </p>
            {currentQ.formula && (
              <div className="p-4 rounded-xl bg-[#09090b] border border-[#1a1a20] font-mono text-center text-sm sm:text-base text-[#38bdf8] tracking-wide">
                {currentQ.formula}
              </div>
            )}
          </div>

          {/* Answer Option Buttons */}
          <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 pt-2">
            {currentQ.options.map((opt) => {
              const isSelected = selectedOption === opt.label;
              const isCorrectOpt = opt.label === currentQ.correctAnswer;
              
              let btnStyle = "bg-[#18181e] text-[#d4d4d8] border-[#2c2c35] hover:border-[#3f3f4e] hover:bg-[#25252e]";
              
              if (selectedOption) {
                if (isCorrectOpt) {
                  btnStyle = "bg-[#10b981]/20 text-white border-[#10b981] shadow-lg shadow-[#10b981]/20";
                } else if (isSelected && !isCorrect) {
                  btnStyle = "bg-[#ef4444]/20 text-white border-[#ef4444]";
                } else {
                  btnStyle = "bg-[#111115] text-[#71717a] border-[#1a1a20] opacity-50";
                }
              }

              return (
                <button
                  key={opt.label}
                  onClick={() => handleSelectOption(opt.label)}
                  disabled={selectedOption !== null}
                  className={`p-4 rounded-xl border text-left font-mono text-xs sm:text-sm transition-all flex items-center justify-between group ${btnStyle}`}
                >
                  <div className="flex items-center gap-3">
                    <span className={`w-7 h-7 rounded-lg flex items-center justify-center font-bold text-xs ${
                      selectedOption && isCorrectOpt
                        ? "bg-[#10b981] text-white"
                        : selectedOption && isSelected && !isCorrect
                        ? "bg-[#ef4444] text-white"
                        : "bg-[#25252e] text-[#94a3b8] group-hover:text-white"
                    }`}>
                      {opt.label}
                    </span>
                    <span className="font-sans font-medium text-white/90">{opt.text}</span>
                  </div>

                  {selectedOption && isCorrectOpt && (
                    <CheckCircle2 className="w-5 h-5 text-[#10b981] shrink-0" />
                  )}
                  {selectedOption && isSelected && !isCorrect && (
                    <XCircle className="w-5 h-5 text-[#ef4444] shrink-0" />
                  )}
                </button>
              );
            })}
          </div>

          {/* Explanation Box (Reveals when an option is clicked) */}
          {showExplanation && (
            <div className={`p-5 rounded-2xl border animate-in fade-in slide-in-from-top-2 duration-300 ${
              isCorrect 
                ? "bg-[#10b981]/10 border-[#10b981]/40" 
                : "bg-[#18181e] border-[#f59e0b]/40"
            }`}>
              <div className="flex items-start gap-3">
                <div className="p-2 rounded-xl bg-black/40 shrink-0 mt-0.5">
                  <Lightbulb className={`w-5 h-5 ${isCorrect ? "text-[#10b981]" : "text-[#f59e0b]"}`} />
                </div>
                <div className="space-y-1.5">
                  <div className="flex items-center gap-2">
                    <span className={`text-xs font-bold ${isCorrect ? "text-[#10b981]" : "text-[#f59e0b]"}`}>
                      {isCorrect ? "🎉 Correct Answer!" : "⚠️ Step-by-Step Solution:"}
                    </span>
                    <span className="text-[10px] text-[#94a3b8]">Correct Option: {currentQ.correctAnswer}</span>
                  </div>
                  <p className="text-xs sm:text-sm text-[#d4d4d8] leading-relaxed">
                    {currentQ.explanation}
                  </p>
                </div>
              </div>

              {/* Next Question / Download Action Bar */}
              <div className="mt-4 pt-4 border-t border-white/10 flex flex-wrap items-center justify-between gap-3">
                <button
                  onClick={handleNextQuestion}
                  className="flex items-center gap-1.5 text-xs font-semibold text-white hover:text-[#10b981] transition-colors"
                >
                  <RefreshCw className="w-3.5 h-3.5" />
                  <span>Try Next Sample Question</span>
                </button>

                <a
                  href="https://play.google.com/store/apps"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-[#10b981] text-black font-bold text-xs hover:bg-[#34d399] transition-all"
                >
                  <Download className="w-3.5 h-3.5" />
                  <span>Practice 10,000+ Questions in App</span>
                </a>
              </div>
            </div>
          )}

        </div>

      </div>
    </section>
  );
}
