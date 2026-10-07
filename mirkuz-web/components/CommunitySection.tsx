"use client";

import React from "react";
import { 
  ShieldCheck, 
  Compass, 
  WifiOff, 
  ArrowRight,
  Sparkles,
  GraduationCap
} from "lucide-react";

export function CommunitySection() {
  const cards = [
    {
      title: "Curriculum Alignment",
      stream: "2017 E.C. Ministry Syllabus",
      description:
        "Laser-focused coverage of the revised Ethiopian national curriculum across all Grade 11 and 12 subjects, eliminating irrelevant topics and focusing purely on matric syllabus goals.",
      icon: ShieldCheck,
      badge: "Ministry Aligned",
    },
    {
      title: "Matric Strategy & Tactics",
      stream: "Pattern Recognition & Speed",
      description:
        "Master recurring entrance examination problem structures, calculus derivations, and speed-solving shortcuts designed to secure competitive 600+ scores for medical & tech faculties.",
      icon: Compass,
      badge: "High-Yield Strategy",
    },
    {
      title: "100% Offline Readiness",
      stream: "Zero Mobile Data Required",
      description:
        "Download full subject video masterclasses, formula sheets, and timed exams directly to your device. Study interruption-free anywhere in Ethiopia without consuming costly internet packages.",
      icon: WifiOff,
      badge: "Offline Ready",
    },
  ];

  return (
    <section id="curriculum" className="bg-[#F8FAFC] py-20 md:py-28 border-b border-[#F1F5F9] relative overflow-hidden">
      {/* Background Soft Glow */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[500px] bg-gradient-to-r from-[#00B894]/5 via-transparent to-[#1E2B58]/5 blur-3xl pointer-events-none" />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        
        {/* Section Header Matching Video Frame 003 */}
        <div className="text-center max-w-2xl mx-auto mb-16">
          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-extrabold text-[#1E2B58] tracking-tight">
            We Focus On
          </h2>
          <div className="title-underline-center" />
          <p className="text-sm sm:text-base text-[#64748B] mt-4">
            A targeted preparation framework designed specifically for Ethiopian matric examinees
          </p>
        </div>

        {/* 3-Card Grid Matching Frame 003 Style */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {cards.map((card, idx) => {
            const Icon = card.icon;
            return (
              <div
                key={idx}
                className="group relative bg-white rounded-3xl p-8 sm:p-9 border border-[#E2E8F0]/80 shadow-dribbble hover:shadow-dribbble-lg transition-all duration-300 hover:-translate-y-1.5 flex flex-col justify-between"
              >
                <div>
                  {/* Icon Container with Mint Background */}
                  <div className="w-14 h-14 rounded-2xl bg-[#E8F8F5] group-hover:bg-[#00B894] flex items-center justify-center transition-colors mb-7">
                    <Icon className="w-7 h-7 text-[#00B894] group-hover:text-white transition-colors" />
                  </div>

                  {/* Title */}
                  <h3 className="text-xl sm:text-2xl font-extrabold text-[#1E2B58] mb-2 leading-snug">
                    {card.title}
                  </h3>

                  {/* Stream Subtitle Badge */}
                  <div className="inline-block px-3 py-1 mb-4 rounded-full bg-[#E8F8F5] text-[#00B894] text-xs font-bold">
                    {card.stream}
                  </div>

                  {/* Description */}
                  <p className="text-sm text-[#64748B] leading-relaxed">
                    {card.description}
                  </p>
                </div>

                {/* Learn More Link */}
                <div className="mt-8 pt-5 border-t border-[#F1F5F9]">
                  <a
                    href="#how-it-works"
                    className="inline-flex items-center gap-2 text-xs font-bold text-[#1E2B58] group-hover:text-[#00B894] transition-colors"
                  >
                    <span>Discover our approach</span>
                    <ArrowRight className="w-3.5 h-3.5 text-[#00B894] group-hover:translate-x-1 transition-transform" />
                  </a>
                </div>
              </div>
            );
          })}
        </div>

      </div>
    </section>
  );
}
