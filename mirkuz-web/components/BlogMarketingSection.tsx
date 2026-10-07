"use client";

import React from "react";
import Image from "next/image";
import { ArrowRight, BookOpen, Clock, Sparkles } from "lucide-react";

export function BlogMarketingSection() {
  const blogs = [
    {
      title: "Top 10 High-Yield Topics in Grade 12 Physics & Calculus for 2017 E.C.",
      tag: "Exam Strategy",
      readTime: "5 min read",
      accentColor: "#4CAF4F",
      bgGradient: "from-emerald-500/20 to-green-600/10",
      svgPattern: "physics",
      link: "https://t.me/mirkuz_exam"
    },
    {
      title: "How to Solve Timed Entrance Aptitude & Logic Questions in Under 45 Seconds",
      tag: "Speed & Accuracy",
      readTime: "4 min read",
      accentColor: "#2684FF",
      bgGradient: "from-blue-500/20 to-sky-600/10",
      svgPattern: "aptitude",
      link: "https://t.me/mirkuz_exam"
    },
    {
      title: "From Regional High School to AAU Medicine: How Kalkidan Scored 642",
      tag: "Success Story",
      readTime: "6 min read",
      accentColor: "#FBC02D",
      bgGradient: "from-amber-500/20 to-yellow-600/10",
      svgPattern: "success",
      link: "https://t.me/mirkuz_exam"
    },
  ];

  return (
    <section id="study-guides" className="bg-[#F8FAFC] py-20 md:py-28 border-b border-[#E2E8F0]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Section Heading */}
        <div className="text-center max-w-2xl mx-auto space-y-3 mb-16">
          <div className="inline-flex items-center gap-2 px-3.5 py-1 rounded-full bg-[#E8F8F5] text-[#00B894] text-xs font-semibold">
            <Sparkles className="w-3.5 h-3.5" />
            <span>Curated Matric Insights</span>
          </div>
          <h2 className="text-3xl sm:text-4xl font-extrabold text-[#1E2B58] tracking-tight">
            National Exam Insights & <span className="title-underline-center text-[#1E2B58]">Study Guides</span>
          </h2>
          <p className="text-sm sm:text-base text-[#64748B] leading-relaxed pt-2">
            Proven strategies, question patterns, and examiner breakdowns to help you conquer the 2017 E.C. matric exam with confidence.
          </p>
        </div>

        {/* 3 Floating-Card Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 lg:gap-8 pb-6">
          {blogs.map((b, idx) => (
            <div key={idx} className="relative flex flex-col items-center group">
              
              {/* Card Image / Visual Illustration Container */}
              <div className={`w-full aspect-16/10 rounded-2xl overflow-hidden bg-gradient-to-br ${b.bgGradient} border border-[#E2E8F0] relative flex items-center justify-center p-6 shadow-xs group-hover:shadow-md transition-all`}>
                
                {/* Visual SVG Graphic representing blog topic */}
                <div className="w-full h-full flex items-center justify-center">
                  {b.svgPattern === "physics" && (
                    <svg viewBox="0 0 200 120" className="w-4/5 h-auto">
                      <circle cx="100" cy="60" r="40" fill="none" stroke="#00B894" strokeWidth="2" strokeDasharray="4 4" />
                      <ellipse cx="100" cy="60" rx="60" ry="25" fill="none" stroke="#00B894" strokeWidth="1.5" transform="rotate(-25 100 60)" />
                      <ellipse cx="100" cy="60" rx="60" ry="25" fill="none" stroke="#1E2B58" strokeWidth="1.5" transform="rotate(25 100 60)" />
                      <circle cx="100" cy="60" r="14" fill="#00B894" />
                      <circle cx="145" cy="45" r="5" fill="#1E2B58" />
                      <circle cx="55" cy="75" r="5" fill="#1E2B58" />
                    </svg>
                  )}

                  {b.svgPattern === "aptitude" && (
                    <svg viewBox="0 0 200 120" className="w-4/5 h-auto">
                      <rect x="40" y="20" width="120" height="80" rx="10" fill="#FFFFFF" stroke="#00B894" strokeWidth="1.5" />
                      <circle cx="65" cy="50" r="14" fill="#00B894" fillOpacity="0.15" />
                      <path d="M60 50 L63 53 L70 47" stroke="#00B894" strokeWidth="2" fill="none" />
                      <rect x="90" y="42" width="55" height="6" rx="3" fill="#1E2B58" />
                      <rect x="90" y="54" width="40" height="5" rx="2.5" fill="#94A3B8" />
                      <line x1="50" y1="78" x2="150" y2="78" stroke="#E2E8F0" strokeWidth="1" />
                      <circle cx="100" cy="88" r="4" fill="#00B894" />
                    </svg>
                  )}

                  {b.svgPattern === "success" && (
                    <svg viewBox="0 0 200 120" className="w-4/5 h-auto">
                      <polygon points="100,15 110,45 142,45 116,64 126,94 100,75 74,94 84,64 58,45 90,45" fill="#F59E0B" />
                      <circle cx="100" cy="55" r="18" fill="#FFFFFF" />
                      <text x="93" y="61" fill="#1E2B58" fontSize="16">🎓</text>
                    </svg>
                  )}
                </div>

                {/* Pill tag top right */}
                <span className="absolute top-3.5 right-3.5 px-3 py-1 rounded-full bg-white/95 backdrop-blur-xs text-[11px] font-bold text-[#1E2B58] shadow-xs">
                  {b.tag}
                </span>
              </div>

              {/* Floating Overlapping White Card */}
              <div className="w-[90%] -mt-10 bg-white rounded-2xl p-6 shadow-dribbble border border-[#E2E8F0] relative z-10 flex flex-col justify-between group-hover:-translate-y-1 transition-transform">
                
                <h3 className="text-base font-bold text-[#1E2B58] leading-snug line-clamp-2 mb-4 group-hover:text-[#00B894] transition-colors">
                  {b.title}
                </h3>

                <div className="flex items-center justify-between pt-2 border-t border-[#F1F5F9]">
                  <a
                    href={b.link}
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1.5 text-xs font-bold text-[#00B894] hover:text-[#009874] hover:gap-2.5 transition-all"
                  >
                    <span>Read on Telegram</span>
                    <ArrowRight className="w-3.5 h-3.5" />
                  </a>

                  <span className="text-[11px] text-[#94A3B8] flex items-center gap-1">
                    <Clock className="w-3 h-3" />
                    <span>{b.readTime}</span>
                  </span>
                </div>

              </div>

            </div>
          ))}
        </div>

      </div>
    </section>
  );
}
