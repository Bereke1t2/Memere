"use client";

import React from "react";

export function QuoteBanner() {
  return (
    <section className="bg-white py-12 sm:py-16">
      <div className="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Glassmorphic / Mint Quote Callout matching Frame 003 of Dribbble Video */}
        <div className="relative rounded-3xl bg-gradient-to-br from-[#E8F8F5]/90 to-[#E8F8F5]/50 border border-[#00B894]/20 p-8 sm:p-14 text-center shadow-dribbble overflow-hidden">
          
          {/* Subtle Ambient Glow */}
          <div className="absolute -top-16 -right-16 w-44 h-44 bg-[#00B894]/15 rounded-full blur-2xl pointer-events-none" />
          <div className="absolute -bottom-16 -left-16 w-44 h-44 bg-[#1E2B58]/10 rounded-full blur-2xl pointer-events-none" />

          {/* Giant Teal Quote Mark (as seen in video) */}
          <div className="text-5xl sm:text-6xl text-[#00B894] font-serif font-black leading-none mb-3 select-none">
            “
          </div>

          {/* Bold Statement */}
          <h3 className="text-2xl sm:text-3xl lg:text-4xl font-extrabold text-[#1E2B58] tracking-tight leading-snug max-w-2xl mx-auto mb-4">
            Not just passing the matric exam. <br className="hidden sm:inline" />
            Securing your top-choice university.
          </h3>

          {/* Subtitle / Philosophy */}
          <p className="text-sm sm:text-base text-[#64748B] max-w-xl mx-auto leading-relaxed">
            When you prepare with real entrance questions and examiner video derivations, guesswork disappears — and test day confidence takes over.
          </p>

        </div>

      </div>
    </section>
  );
}
