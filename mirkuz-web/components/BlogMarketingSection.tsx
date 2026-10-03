"use client";

import React from "react";
import Image from "next/image";
import { ArrowRight, BookOpen, Clock, Sparkles } from "lucide-react";

export function BlogMarketingSection() {
  const blogs = [
    {
      title: "Creating Streamlined Safeguarding Processes with OneRen",
      sub: "Top 10 High-Yield Topics in Grade 12 Physics & Calculus for 2017 E.C.",
      tag: "Exam Strategy",
      readTime: "5 min read",
      accentColor: "#4CAF4F",
      bgGradient: "from-emerald-500/20 to-green-600/10",
      svgPattern: "physics"
    },
    {
      title: "What are your safeguarding responsibilities and how can you manage them?",
      sub: "How to Solve Timed Entrance Aptitude & Logic Questions in Under 45 Seconds",
      tag: "Speed & Accuracy",
      readTime: "4 min read",
      accentColor: "#2684FF",
      bgGradient: "from-blue-500/20 to-sky-600/10",
      svgPattern: "aptitude"
    },
    {
      title: "Revamping the Membership Model with Triathlon Australia",
      sub: "From Hawassa to AAU Medicine: How Kalkidan Scored 642 Using Offline Prep",
      tag: "Success Case Study",
      readTime: "6 min read",
      accentColor: "#FBC02D",
      bgGradient: "from-amber-500/20 to-yellow-600/10",
      svgPattern: "success"
    },
  ];

  return (
    <section className="bg-white py-20 md:py-28 border-b border-[#E8ECF2]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Section Heading */}
        <div className="text-center max-w-2xl mx-auto space-y-3 mb-16">
          <h2 className="text-3xl sm:text-4xl font-bold text-[#4D4D4D] tracking-tight">
            Caring is the new marketing
          </h2>
          <p className="text-sm sm:text-base text-[#717171] leading-relaxed">
            The Nextcent blog is the best place to read about the latest membership insights, trends and more. See who’s joining the community, read about how our community are increasing their membership income and lot’s more.
          </p>
        </div>

        {/* 3 Floating-Card Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-10 lg:gap-8 pb-10">
          {blogs.map((b, idx) => (
            <div key={idx} className="relative flex flex-col items-center group">
              
              {/* Card Image / Visual Illustration Container */}
              <div className={`w-full aspect-16/10 rounded-xl overflow-hidden bg-gradient-to-br ${b.bgGradient} border border-[#E8ECF2] relative flex items-center justify-center p-6 shadow-xs group-hover:shadow-md transition-all`}>
                
                {/* Visual SVG Graphic representing blog topic */}
                <div className="w-full h-full flex items-center justify-center">
                  {b.svgPattern === "physics" && (
                    <svg viewBox="0 0 200 120" className="w-4/5 h-auto">
                      <circle cx="100" cy="60" r="40" fill="none" stroke="#4CAF4F" strokeWidth="2" strokeDasharray="4 4" />
                      <ellipse cx="100" cy="60" rx="60" ry="25" fill="none" stroke="#4CAF4F" strokeWidth="1.5" transform="rotate(-25 100 60)" />
                      <ellipse cx="100" cy="60" rx="60" ry="25" fill="none" stroke="#4CAF4F" strokeWidth="1.5" transform="rotate(25 100 60)" />
                      <circle cx="100" cy="60" r="14" fill="#4CAF4F" />
                      <circle cx="145" cy="45" r="5" fill="#388E3C" />
                      <circle cx="55" cy="75" r="5" fill="#388E3C" />
                    </svg>
                  )}

                  {b.svgPattern === "aptitude" && (
                    <svg viewBox="0 0 200 120" className="w-4/5 h-auto">
                      <rect x="40" y="20" width="120" height="80" rx="8" fill="#FFFFFF" stroke="#2684FF" strokeWidth="1.5" />
                      <circle cx="65" cy="50" r="14" fill="#2684FF" fillOpacity="0.15" />
                      <path d="M60 50 L63 53 L70 47" stroke="#2684FF" strokeWidth="2" fill="none" />
                      <rect x="90" y="42" width="55" height="6" rx="3" fill="#2684FF" />
                      <rect x="90" y="54" width="40" height="5" rx="2.5" fill="#ABBED1" />
                      <line x1="50" y1="78" x2="150" y2="78" stroke="#E8ECF2" strokeWidth="1" />
                      <circle cx="100" cy="88" r="4" fill="#2684FF" />
                    </svg>
                  )}

                  {b.svgPattern === "success" && (
                    <svg viewBox="0 0 200 120" className="w-4/5 h-auto">
                      <polygon points="100,15 110,45 142,45 116,64 126,94 100,75 74,94 84,64 58,45 90,45" fill="#FBC02D" />
                      <circle cx="100" cy="55" r="18" fill="#FFFFFF" />
                      <text x="94" y="60" fill="#263238" fontSize="14" fontWeight="bold">🎓</text>
                    </svg>
                  )}
                </div>

                {/* Pill tag top right */}
                <span className="absolute top-3 right-3 px-2.5 py-1 rounded-full bg-white/90 backdrop-blur-xs text-[11px] font-semibold text-[#4D4D4D] shadow-xs">
                  {b.tag}
                </span>
              </div>

              {/* Floating Overlapping White Card (Signature Nexcent Element) */}
              <div className="w-[88%] -mt-12 bg-white rounded-lg p-5 sm:p-6 shadow-[0_8px_16px_rgba(171,190,209,0.4)] border border-[#E8ECF2] relative z-10 flex flex-col justify-between group-hover:-translate-y-1 transition-transform">
                
                <h3 className="text-base sm:text-lg font-bold text-[#4D4D4D] leading-snug line-clamp-2 mb-4">
                  {b.title}
                </h3>

                <div className="flex items-center justify-between pt-2">
                  <a
                    href="#blog"
                    className="inline-flex items-center gap-1.5 text-sm font-bold text-[#4CAF4F] hover:text-[#388E3C] hover:gap-2.5 transition-all"
                  >
                    <span>Readmore</span>
                    <ArrowRight className="w-4 h-4" />
                  </a>

                  <span className="text-[11px] text-[#89939E] flex items-center gap-1">
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
