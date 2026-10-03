"use client";

import React from "react";
import Image from "next/image";
import { ArrowRight, CheckCircle2, Calendar, FileText, WifiOff, Zap } from "lucide-react";

export function FeatureSplitTwo() {
  return (
    <section id="product" className="bg-white py-20 md:py-28 border-b border-[#E8ECF2] overflow-hidden">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-16 items-center">
          
          {/* Left Column: Copy & Button */}
          <div className="lg:col-span-7 flex flex-col items-start space-y-6 order-2 lg:order-1">
            
            <div className="space-y-3">
              <h2 className="text-3xl sm:text-4xl font-bold text-[#4D4D4D] tracking-tight leading-tight">
                How to design your site footer <br className="hidden sm:inline" />
                like we did
              </h2>

              <div className="inline-block px-3 py-1 rounded-md bg-[#E8F5E9] text-[#4CAF4F] text-xs font-semibold">
                High-Yield Timed Mock Exam Engine
              </div>
            </div>

            <p className="text-sm sm:text-base text-[#717171] leading-relaxed max-w-xl">
              Donec tempus, odio eget scelerisque luctus, feugiat sem neque pellentesque ipsum, a porta nisi odio ac purus. Fusce feugiat dui lorem, a vestibulum magna finibus sed. Donec tempus, odio eget scelerisque luctus, feugiat sem neque pellentesque ipsum, a porta nisi odio ac purus. Fusce feugiat dui lorem, a vestibulum magna finibus sed.
            </p>

            {/* Pillar list */}
            <div className="space-y-2.5 w-full text-sm text-[#4D4D4D]">
              <div className="flex items-start gap-2.5">
                <CheckCircle2 className="w-4 h-4 text-[#4CAF4F] mt-0.5 shrink-0" />
                <span>Real 2-hour countdown timer identical to Ministry of Education matrix exams</span>
              </div>
              <div className="flex items-start gap-2.5">
                <CheckCircle2 className="w-4 h-4 text-[#4CAF4F] mt-0.5 shrink-0" />
                <span>60-question jump palette with flag-for-review bookmarks</span>
              </div>
              <div className="flex items-start gap-2.5">
                <CheckCircle2 className="w-4 h-4 text-[#4CAF4F] mt-0.5 shrink-0" />
                <span>Immediate answer explanations with formula breakdowns and step-by-step videos</span>
              </div>
            </div>

            {/* Learn More Button */}
            <a
              href="#exam-demo"
              className="inline-flex items-center gap-2 px-8 py-3 rounded-md bg-[#4CAF4F] hover:bg-[#388E3C] text-white font-medium text-sm sm:text-base shadow-sm hover:shadow transition-all hover:-translate-y-0.5 active:translate-y-0"
            >
              <span>Learn More</span>
              <ArrowRight className="w-4 h-4" />
            </a>

          </div>

          {/* Right Column: Clean Vector Mobile & Calendar Illustration */}
          <div className="lg:col-span-5 flex justify-center order-1 lg:order-2">
            <div className="relative w-full max-w-[420px] aspect-square flex items-center justify-center">
              
              <svg viewBox="0 0 400 400" className="w-full h-auto drop-shadow-sm">
                <defs>
                  <linearGradient id="circleGrad2" x1="0" y1="0" x2="1" y2="1">
                    <stop offset="0%" stopColor="#F5F7FA" />
                    <stop offset="100%" stopColor="#E8F5E9" />
                  </linearGradient>
                </defs>

                {/* Background soft circle */}
                <circle cx="200" cy="200" r="160" fill="url(#circleGrad2)" />

                {/* Calendar / Schedule Card */}
                <rect x="70" y="80" width="260" height="240" rx="14" fill="#FFFFFF" filter="drop-shadow(0 6px 16px rgba(0,0,0,0.06))" />
                
                {/* Header bar */}
                <rect x="70" y="80" width="260" height="48" rx="14" fill="#4CAF4F" />
                <text x="94" y="110" fill="#FFFFFF" fontSize="13" fontWeight="bold">Exam Preparation Roadmap</text>
                <circle cx="295" cy="104" r="8" fill="#FFFFFF" fillOpacity="0.25" />
                <path d="M292 104 L294 106 L298 102" stroke="#FFFFFF" strokeWidth="1.5" fill="none" />

                {/* Calendar Rows */}
                {/* Item 1 */}
                <rect x="90" y="145" width="220" height="40" rx="6" fill="#F5F7FA" />
                <circle cx="110" cy="165" r="10" fill="#E8F5E9" />
                <path d="M106 165 L109 168 L114 163" stroke="#4CAF4F" strokeWidth="2" fill="none" />
                <text x="130" y="160" fill="#4D4D4D" fontSize="11" fontWeight="bold">Mathematics (Calculus)</text>
                <text x="130" y="173" fill="#717171" fontSize="9">2016 E.C. Exam Completed — 94%</text>

                {/* Item 2 */}
                <rect x="90" y="195" width="220" height="40" rx="6" fill="#F5F7FA" />
                <circle cx="110" cy="215" r="10" fill="#E8F5E9" />
                <path d="M106 215 L109 218 L114 213" stroke="#4CAF4F" strokeWidth="2" fill="none" />
                <text x="130" y="210" fill="#4D4D4D" fontSize="11" fontWeight="bold">Physics (Electromagnetism)</text>
                <text x="130" y="223" fill="#717171" fontSize="9">Unit 4 & 5 Video Mastery — 100%</text>

                {/* Item 3 */}
                <rect x="90" y="245" width="220" height="40" rx="6" fill="#F5F7FA" />
                <circle cx="110" cy="265" r="10" fill="#E8F5E9" />
                <circle cx="110" cy="265" r="4" fill="#4CAF4F" />
                <text x="130" y="260" fill="#4D4D4D" fontSize="11" fontWeight="bold">Chemistry & Aptitude</text>
                <text x="130" y="273" fill="#717171" fontSize="9">Timed Mock Simulation in Progress</text>

                {/* Floating Timer Badge */}
                <g className="animate-pulse">
                  <rect x="40" y="230" width="95" height="42" rx="8" fill="#FFFFFF" filter="drop-shadow(0 4px 10px rgba(0,0,0,0.1))" />
                  <circle cx="58" cy="251" r="10" fill="#E8F5E9" />
                  <text x="54" y="255" fill="#4CAF4F" fontSize="10" fontWeight="bold">⏱</text>
                  <text x="76" y="247" fill="#4D4D4D" fontSize="10" fontWeight="bold">01:45:00</text>
                  <text x="76" y="258" fill="#717171" fontSize="8">Remaining Time</text>
                </g>
              </svg>

            </div>
          </div>

        </div>
      </div>
    </section>
  );
}
