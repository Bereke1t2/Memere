"use client";

import React from "react";
import { 
  BookOpen, 
  CheckSquare, 
  Clock, 
  BarChart3, 
  Sparkles,
  ArrowRight
} from "lucide-react";

export function HowWeDoIt() {
  const steps = [
    {
      num: "01",
      title: "Concept Mastery",
      description:
        "Watch concise unit breakdowns taught by top Ethiopian educators. Understand difficult formulas, proofs, and principles without rote memorization.",
      icon: BookOpen,
    },
    {
      num: "02",
      title: "Past Exam Drills",
      description:
        "Solve 10,000+ verified national entrance questions from 2008 to 2016 E.C. Learn recurring exam question traps and examiner answer strategies.",
      icon: CheckSquare,
    },
    {
      num: "03",
      title: "Timed Mock Tests",
      description:
        "Experience the official 2-hour Ministry of Education examination pacing. Train your speed with 60-question jump palettes and bookmarks.",
      icon: Clock,
    },
    {
      num: "04",
      title: "Diagnostic Analytics",
      description:
        "Receive instant AI diagnosis of your weak chapters and formulas before exam day. Track your national score rank and AAU/ASTU cutoff chance.",
      icon: BarChart3,
    },
  ];

  return (
    <section id="how-it-works" className="bg-white py-20 md:py-28 border-b border-[#F1F5F9] relative overflow-hidden">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Section Header with Underline */}
        <div className="text-center max-w-3xl mx-auto mb-20">
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-extrabold text-[#1E2B58] tracking-tight">
            How we do it
          </h2>
          <div className="title-underline-center" />
          <p className="text-xs sm:text-sm text-[#64748B] mt-4 max-w-xl mx-auto leading-relaxed">
            Everything designed to avoid guesswork when it comes to concrete study actions that will elevate your national matric rank.
          </p>
        </div>

        {/* 4-Step Horizontal Timeline Track Matching Frame 004/005 */}
        <div className="relative mb-20">
          
          {/* Dotted Connecting Horizontal Line on Desktop */}
          <div className="hidden lg:block absolute top-7 left-[8%] right-[8%] h-0.5 border-t-2 border-dashed border-[#CBD5E1] z-0" />

          {/* 4 Cards Grid */}
          <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-8 relative z-10">
            {steps.map((step, idx) => {
              const Icon = step.icon;
              return (
                <div key={idx} className="flex flex-col items-center text-center group">
                  
                  {/* Step Number Circle (as seen in Dribbble video) */}
                  <div className="w-14 h-14 rounded-full bg-white border-2 border-[#00B894] flex items-center justify-center font-extrabold text-sm text-[#00B894] shadow-md shadow-[#00B894]/15 mb-6 group-hover:scale-110 group-hover:bg-[#E8F8F5] transition-all">
                    {step.num}
                  </div>

                  {/* Card Container */}
                  <div className="w-full bg-white rounded-3xl p-7 border border-[#E2E8F0] shadow-dribbble hover:shadow-dribbble-lg transition-all duration-300 hover:-translate-y-1.5 flex-1 flex flex-col justify-between">
                    <div>
                      {/* Icon */}
                      <div className="w-12 h-12 rounded-2xl bg-[#E8F8F5] group-hover:bg-[#00B894] mx-auto flex items-center justify-center text-[#00B894] group-hover:text-white transition-colors mb-5">
                        <Icon className="w-6 h-6" />
                      </div>

                      {/* Title */}
                      <h3 className="text-lg font-extrabold text-[#1E2B58] mb-3">
                        {step.title}
                      </h3>

                      {/* Description */}
                      <p className="text-xs sm:text-sm text-[#64748B] leading-relaxed">
                        {step.description}
                      </p>
                    </div>
                  </div>

                </div>
              );
            })}
          </div>

        </div>

        {/* Deep Navy Banner 1 Matching Frame 004/005 */}
        <div className="rounded-3xl bg-[#1E2B58] p-8 sm:p-14 text-center text-white shadow-2xl shadow-[#1E2B58]/20 relative overflow-hidden">
          {/* Subtle Ambient Radial Glow */}
          <div className="absolute top-0 right-0 w-80 h-80 bg-[#00B894]/20 rounded-full blur-3xl pointer-events-none" />
          <div className="absolute bottom-0 left-0 w-60 h-60 bg-[#2A3B72]/40 rounded-full blur-2xl pointer-events-none" />

          <div className="relative z-10 max-w-2xl mx-auto space-y-3">
            <span className="text-xs font-extrabold tracking-widest text-[#00B894] uppercase">
              NO MORE GUESS WORK
            </span>
            <h3 className="text-2xl sm:text-4xl font-extrabold tracking-tight leading-tight text-white">
              Start Working Creatively and Based on Data
            </h3>
            <p className="text-xs sm:text-sm text-[#CBD5E1] pt-1">
              Join 50,000+ Grade 12 students taking control of their university admission future with Mirkuz.
            </p>
          </div>
        </div>

      </div>
    </section>
  );
}
