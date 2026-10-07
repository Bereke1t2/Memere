"use client";

import React from "react";
import { 
  Laptop, 
  RotateCcw, 
  BrainCircuit, 
  Crosshair, 
  ArrowRight,
  Sparkles
} from "lucide-react";

export function FeatureSplitTwo() {
  const cards = [
    {
      title: "Adapting to Digital Exams",
      desc: "Train directly on the Ministry's 2-hour digital testing environment with timed pacing, flag-for-review, and zero negative-marking surprises.",
      icon: Laptop,
    },
    {
      title: "Updating Chapter Mastery",
      desc: "Systematic unit checklists across all 6 subjects ensure every high-frequency topic from Grade 11 & 12 is thoroughly mastered.",
      icon: RotateCcw,
    },
    {
      title: "Long-Term Retention",
      desc: "Downloadable formula sheets, chapter summaries, and spaced-repetition drills solidify key concepts in your memory before exam week.",
      icon: BrainCircuit,
    },
    {
      title: "Targeted Score Recovery",
      desc: "Instant error analysis diagnoses why an answer was wrong and provides immediate step-by-step video derivations.",
      icon: Crosshair,
    },
  ];

  return (
    <section id="features" className="bg-[#F8FAFC] py-20 md:py-28 border-b border-[#F1F5F9] relative overflow-hidden">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-14 items-center">
          
          {/* Left Column: Heading with Underline + Modern Isometric Vector Art */}
          <div className="lg:col-span-5 flex flex-col items-start space-y-8">
            <div>
              <h2 className="text-3xl sm:text-4xl lg:text-5xl font-extrabold text-[#1E2B58] tracking-tight leading-tight">
                The best time to <br />
                prepare with us
              </h2>
              <div className="title-underline" />
              <p className="text-sm text-[#64748B] mt-4 leading-relaxed">
                Whether you have 6 months or 6 weeks before the national matric, Mirkuz accelerates your score trajectory with proven diagnostic workflows.
              </p>
            </div>

            {/* Isometric / 3D Dashboard Art matching Frame 006 */}
            <div className="w-full max-w-[420px] aspect-4/3 rounded-3xl bg-white border border-[#E2E8F0] shadow-dribbble p-5 flex items-center justify-center relative overflow-hidden">
              <svg viewBox="0 0 380 280" className="w-full h-auto drop-shadow-sm">
                <defs>
                  <linearGradient id="isoGrad" x1="0" y1="0" x2="1" y2="1">
                    <stop offset="0%" stopColor="#1E2B58" />
                    <stop offset="100%" stopColor="#162145" />
                  </linearGradient>
                </defs>

                {/* Isometric Board Base */}
                <polygon points="190,40 340,110 190,180 40,110" fill="#E8F8F5" />
                <polygon points="40,110 190,180 190,195 40,125" fill="#CBD5E1" />
                <polygon points="340,110 190,180 190,195 340,125" fill="#94A3B8" />

                {/* Isometric Dashboard Screen Standing on the board */}
                <polygon points="90,140 250,60 250,180 90,260" fill="url(#isoGrad)" />
                {/* Screen top border */}
                <polygon points="90,140 250,60 255,62 95,142" fill="#00B894" />

                {/* Dashboard Elements inside screen */}
                <line x1="110" y1="160" x2="230" y2="100" stroke="#00B894" strokeWidth="3" />
                <line x1="110" y1="180" x2="190" y2="140" stroke="#64748B" strokeWidth="2" />
                <line x1="110" y1="200" x2="210" y2="150" stroke="#64748B" strokeWidth="2" />
                <circle cx="170" cy="130" r="16" fill="#00B894" opacity="0.3" />
                <circle cx="170" cy="130" r="8" fill="#00B894" />

                {/* Isometric Bar Graph on right of the board */}
                <polygon points="270,120 290,110 290,70 270,80" fill="#00B894" />
                <polygon points="290,110 305,117 305,77 290,70" fill="#00A381" />
                <polygon points="270,80 290,70 305,77 285,87" fill="#E8F8F5" />

                {/* Small character interacting with screen (from video) */}
                <circle cx="130" cy="225" r="9" fill="#F8CBA6" />
                <polygon points="124,235 136,235 140,270 120,270" fill="#1E2B58" />
                <line x1="130" y1="245" x2="155" y2="215" stroke="#F8CBA6" strokeWidth="4" strokeLinecap="round" />
              </svg>

              {/* Floating Live Badge */}
              <div className="absolute bottom-3 right-3 px-3 py-1.5 rounded-xl bg-[#1E2B58] text-white text-[11px] font-bold shadow-md flex items-center gap-1.5">
                <span className="w-2 h-2 rounded-full bg-[#00B894] animate-pulse" />
                <span>Live Analytics</span>
              </div>
            </div>
          </div>

          {/* Right Column: 2x2 Grid of Feature Cards Matching Frame 006 */}
          <div className="lg:col-span-7 grid grid-cols-1 sm:grid-cols-2 gap-6">
            {cards.map((c, i) => {
              const Icon = c.icon;
              return (
                <div
                  key={i}
                  className="bg-white rounded-3xl p-7 border border-[#E2E8F0] shadow-dribbble hover:shadow-dribbble-lg transition-all duration-300 hover:-translate-y-1 group flex flex-col justify-between"
                >
                  <div>
                    {/* Icon Container with Mint Background */}
                    <div className="w-13 h-13 rounded-2xl bg-[#E8F8F5] group-hover:bg-[#00B894] flex items-center justify-center text-[#00B894] group-hover:text-white transition-colors mb-5">
                      <Icon className="w-6 h-6" />
                    </div>

                    {/* Title */}
                    <h3 className="text-lg font-extrabold text-[#1E2B58] mb-2 leading-snug">
                      {c.title}
                    </h3>

                    {/* Description */}
                    <p className="text-xs sm:text-sm text-[#64748B] leading-relaxed">
                      {c.desc}
                    </p>
                  </div>
                </div>
              );
            })}
          </div>

        </div>

      </div>
    </section>
  );
}
