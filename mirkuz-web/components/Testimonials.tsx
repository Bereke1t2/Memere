"use client";

import React from "react";
import { Star, Quote, Award, GraduationCap, CheckCircle2 } from "lucide-react";

export function Testimonials() {
  const reviews = [
    {
      name: "Yared Tadesse",
      stream: "Natural Science (2016 E.C.)",
      score: "612 / 700",
      placement: "Addis Ababa University — School of Medicine",
      quote:
        "Mirkuz completely changed how I revised Physics and Math. The timed mock exams trained my pacing so well that on actual exam day, I finished 15 minutes early without panic.",
      location: "Addis Ababa",
    },
    {
      name: "Selamawit Bekele",
      stream: "Social Science (2016 E.C.)",
      score: "564 / 600",
      placement: "Addis Ababa University — Law School",
      quote:
        "The offline download feature was a lifesaver. I live in Hawassa and whenever internet was spotty, all my History and Economics video lessons were already saved on my phone.",
      location: "Hawassa",
    },
    {
      name: "Abel Getachew",
      stream: "Natural Science (2015 E.C.)",
      score: "598 / 700",
      placement: "ASTU — Electrical & Computer Engineering",
      quote:
        "The step-by-step video solutions for past entrance questions from 2008 to 2015 E.C. made difficult calculus questions feel simple. Best investment for Grade 12.",
      location: "Bahir Dar",
    },
  ];

  return (
    <section className="py-24 bg-[#09090d] relative overflow-hidden border-t border-[#1a1a20]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        
        {/* Section Header */}
        <div className="text-center max-w-3xl mx-auto mb-16 space-y-4">
          <div className="inline-flex items-center gap-2 px-3.5 py-1 rounded-full bg-[#f59e0b]/15 text-[#f59e0b] border border-[#f59e0b]/30 text-xs font-bold">
            <GraduationCap className="w-3.5 h-3.5" />
            <span>Proven Student Outcomes</span>
          </div>
          <h2 className="text-3xl sm:text-4xl font-extrabold text-white tracking-tight">
            Trusted by Ethiopia’s Top University Entrance Scorers
          </h2>
          <p className="text-sm sm:text-base text-[#94a3b8]">
            Hear from students across Ethiopia who prepared with Mirkuz and earned placements at top universities.
          </p>
        </div>

        {/* Testimonials Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
          {reviews.map((rev, idx) => (
            <div
              key={idx}
              className="p-7 rounded-3xl bg-[#111115] border border-[#1a1a20] hover:border-[#2c2c35] transition-all flex flex-col justify-between space-y-6 group"
            >
              <div className="space-y-4">
                {/* Rating & Score Badge */}
                <div className="flex items-center justify-between">
                  <div className="flex text-[#f59e0b]">
                    {[...Array(5)].map((_, i) => (
                      <Star key={i} className="w-4 h-4 fill-current" />
                    ))}
                  </div>
                  <span className="px-2.5 py-1 rounded-full bg-[#10b981]/15 text-[#10b981] font-mono text-xs font-bold border border-[#10b981]/30">
                    Entrance Score: {rev.score}
                  </span>
                </div>

                {/* Quote Text */}
                <p className="text-xs sm:text-sm text-[#d4d4d8] leading-relaxed italic">
                  &ldquo;{rev.quote}&rdquo;
                </p>
              </div>

              {/* Student Metadata */}
              <div className="pt-4 border-t border-[#1a1a20] space-y-1">
                <div className="flex items-center justify-between">
                  <h3 className="text-sm font-bold text-white">{rev.name}</h3>
                  <span className="text-[11px] text-[#71717a]">{rev.location}</span>
                </div>
                <div className="text-xs text-[#10b981] font-medium">{rev.placement}</div>
                <div className="text-[11px] text-[#94a3b8]">{rev.stream}</div>
              </div>
            </div>
          ))}
        </div>

      </div>
    </section>
  );
}
