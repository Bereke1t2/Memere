"use client";

import React from "react";
import { ArrowRight, CheckCircle2, Sparkles, BookOpen, Clock, Zap } from "lucide-react";

export function FeatureSplitOne() {
  return (
    <section id="features" className="bg-white py-20 md:py-28 border-b border-[#E8ECF2] overflow-hidden">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-16 items-center">
          
          {/* Left Column: Clean Vector Illustration */}
          <div className="lg:col-span-5 flex justify-center">
            <div className="relative w-full max-w-[420px] aspect-square flex items-center justify-center">
              
              <svg viewBox="0 0 400 400" className="w-full h-auto drop-shadow-sm">
                <defs>
                  <linearGradient id="circleGrad" x1="0" y1="0" x2="1" y2="1">
                    <stop offset="0%" stopColor="#E8F5E9" />
                    <stop offset="100%" stopColor="#F5F7FA" />
                  </linearGradient>
                </defs>

                {/* Soft background circle */}
                <circle cx="200" cy="200" r="160" fill="url(#circleGrad)" />

                {/* Mobile / Tablet Screen in center */}
                <rect x="130" y="80" width="160" height="250" rx="16" fill="#263238" />
                <rect x="138" y="90" width="144" height="230" rx="10" fill="#FFFFFF" />
                
                {/* Screen Header */}
                <rect x="148" y="102" width="60" height="8" rx="4" fill="#4CAF4F" />
                <circle cx="266" cy="106" r="5" fill="#E8F5E9" />
                <circle cx="266" cy="106" r="2.5" fill="#4CAF4F" />

                {/* Video Lesson Preview inside Screen */}
                <rect x="148" y="120" width="124" height="65" rx="6" fill="#F5F7FA" />
                <circle cx="210" cy="152" r="14" fill="#4CAF4F" />
                <polygon points="206,146 218,152 206,158" fill="#FFFFFF" />

                {/* Chapter List inside Screen */}
                <rect x="148" y="196" width="124" height="26" rx="4" fill="#E8F5E9" />
                <circle cx="160" cy="209" r="5" fill="#4CAF4F" />
                <rect x="172" y="206" width="70" height="6" rx="3" fill="#263238" />

                <rect x="148" y="230" width="124" height="26" rx="4" fill="#F5F7FA" />
                <circle cx="160" cy="243" r="5" fill="#ABBED1" />
                <rect x="172" y="240" width="80" height="6" rx="3" fill="#717171" />

                <rect x="148" y="264" width="124" height="26" rx="4" fill="#F5F7FA" />
                <circle cx="160" cy="277" r="5" fill="#ABBED1" />
                <rect x="172" y="274" width="60" height="6" rx="3" fill="#717171" />

                {/* Floating Achievement Badge Top Left */}
                <g className="animate-bounce" style={{ animationDuration: "4s" }}>
                  <rect x="50" y="120" width="100" height="42" rx="8" fill="#FFFFFF" filter="drop-shadow(0 4px 8px rgba(0,0,0,0.06))" />
                  <circle cx="70" cy="141" r="10" fill="#E8F5E9" />
                  <path d="M66 141 L69 144 L75 138" stroke="#4CAF4F" strokeWidth="2" fill="none" />
                  <text x="86" y="137" fill="#4D4D4D" fontSize="10" fontWeight="bold">100% Passed</text>
                  <text x="86" y="148" fill="#717171" fontSize="8">National Standard</text>
                </g>

                {/* Floating Metrics Badge Bottom Right */}
                <g>
                  <rect x="250" y="260" width="110" height="48" rx="8" fill="#FFFFFF" filter="drop-shadow(0 4px 8px rgba(0,0,0,0.06))" />
                  <circle cx="272" cy="284" r="12" fill="#E8F5E9" />
                  <text x="268" y="288" fill="#4CAF4F" fontSize="12" fontWeight="bold">★</text>
                  <text x="292" y="280" fill="#4D4D4D" fontSize="10" fontWeight="bold">600+ Score</text>
                  <text x="292" y="292" fill="#717171" fontSize="8">Top University Entry</text>
                </g>
              </svg>

            </div>
          </div>

          {/* Right Column: Copy & Button */}
          <div className="lg:col-span-7 flex flex-col items-start space-y-6">
            
            <div className="space-y-3">
              <h2 className="text-3xl sm:text-4xl font-bold text-[#4D4D4D] tracking-tight leading-tight">
                The unseen of spending three <br className="hidden sm:inline" />
                years at Pixelgrade
              </h2>
              
              <div className="inline-block px-3 py-1 rounded-md bg-[#E8F5E9] text-[#4CAF4F] text-xs font-semibold">
                Ethiopian National Entrance Exam Preparation
              </div>
            </div>

            <p className="text-sm sm:text-base text-[#717171] leading-relaxed max-w-xl">
              Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed sit amet justo ipsum. Sed accumsan quam vitae est varius fringilla. Pellentesque placerat vestibulum lorem sed porta. Nullam mattis tristique iaculis. Nullam pulvinar sit amet risus pretium auctor. Etiam quis massa pulvinar, aliquam quam vitae, tempus sem. Donec elementum pulvinar odio.
            </p>

            {/* Benefit Highlights */}
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-3 w-full py-1 text-sm text-[#4D4D4D]">
              <div className="flex items-center gap-2">
                <CheckCircle2 className="w-4 h-4 text-[#4CAF4F] shrink-0" />
                <span>2008 – 2016 E.C. Matric Exam Solutions</span>
              </div>
              <div className="flex items-center gap-2">
                <CheckCircle2 className="w-4 h-4 text-[#4CAF4F] shrink-0" />
                <span>Chapter-by-Chapter Video Courses</span>
              </div>
              <div className="flex items-center gap-2">
                <CheckCircle2 className="w-4 h-4 text-[#4CAF4F] shrink-0" />
                <span>100% Encrypted Offline Phone Storage</span>
              </div>
              <div className="flex items-center gap-2">
                <CheckCircle2 className="w-4 h-4 text-[#4CAF4F] shrink-0" />
                <span>Instant Telebirr & Chapa Activation</span>
              </div>
            </div>

            {/* Green Action CTA */}
            <a
              href="https://play.google.com/store/apps"
              target="_blank"
              rel="noopener noreferrer"
              className="inline-flex items-center gap-2 px-8 py-3 rounded-md bg-[#4CAF4F] hover:bg-[#388E3C] text-white font-medium text-sm sm:text-base shadow-sm hover:shadow transition-all hover:-translate-y-0.5 active:translate-y-0"
            >
              <span>Learn More</span>
              <ArrowRight className="w-4 h-4" />
            </a>

          </div>

        </div>
      </div>
    </section>
  );
}
