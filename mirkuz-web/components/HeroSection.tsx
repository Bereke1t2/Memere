"use client";

import React, { useState } from "react";
import { 
  ArrowRight, 
  CheckCircle2, 
  Sparkles, 
  Smartphone, 
  Send,
  BookOpen,
  Award,
  Zap
} from "lucide-react";

export function HeroSection() {
  return (
    <section id="hero" className="relative pt-32 pb-20 md:pt-44 md:pb-28 overflow-hidden bg-white">
      {/* Soft Ambient Radial Background Mesh (from Dribbble video) */}
      <div className="absolute top-0 right-0 w-[650px] h-[650px] bg-gradient-to-br from-[#00B894]/8 via-[#00B894]/2 to-transparent rounded-full blur-3xl pointer-events-none -mr-40 -mt-20" />
      <div className="absolute top-1/4 left-0 w-[550px] h-[550px] bg-gradient-to-tr from-[#1E2B58]/6 via-transparent to-transparent rounded-full blur-3xl pointer-events-none -ml-40" />

      {/* Floating Confetti / Gem Dots (Signature Dribbble Video detail) */}
      <div className="absolute top-36 left-1/2 w-2 h-2 rounded-full bg-[#00B894] opacity-75 animate-ping hidden sm:block" />
      <div className="absolute top-48 right-1/4 w-2.5 h-2.5 rounded-sm bg-[#1E2B58] opacity-60 rotate-45 hidden sm:block" />
      <div className="absolute bottom-28 left-1/4 w-3 h-3 rounded-full bg-[#00B894]/30 hidden sm:block" />
      <div className="absolute bottom-40 right-10 w-2 h-2 rounded-sm bg-[#00B894] rotate-12 hidden sm:block" />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-8 items-center min-h-[500px]">
          
          {/* Left Column: Typography & Action Buttons */}
          <div className="lg:col-span-6 flex flex-col items-start space-y-7">
            

            {/* Main Two-Tone Headline (Exact Dribbble Style) */}
            <h1 className="text-4xl sm:text-6xl lg:text-[68px] font-extrabold tracking-tight text-[#1E2B58] leading-[1.12]">
              Study smart <br />
              Score <span className="text-[#00B894]">higher</span>
            </h1>

            {/* Subtitle */}
            <div className="space-y-2">
              <h2 className="text-2xl sm:text-3xl font-bold text-[#1E2B58]/90 tracking-tight">
                Steering students to <br className="hidden sm:inline" />
                achieve university admission
              </h2>
              <p className="text-base sm:text-lg text-[#64748B] max-w-lg leading-relaxed pt-1">
                Master Ethiopia's Grade 12 National Entrance Exam with senior examiner video masterclasses, 10,000+ solved past matric questions, and 100% offline study.
              </p>
            </div>

            {/* Action Buttons Matching Video: Deep Navy Pill Button */}
            <div className="flex flex-wrap items-center gap-4 pt-2">
              <a
                href="https://play.google.com/store/apps/details?id=et.mirkuz.mobile"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center justify-center gap-2.5 px-8 py-4 rounded-xl bg-[#1E2B58] hover:bg-[#162145] text-white font-bold text-sm sm:text-base shadow-xl shadow-[#1E2B58]/25 hover:shadow-2xl hover:-translate-y-0.5 active:translate-y-0 transition-all"
              >
                <Smartphone className="w-4 h-4 text-[#00B894]" />
                <span>Download App Free</span>
                <ArrowRight className="w-4 h-4" />
              </a>

              <a
                href="https://t.me/mirkuz_exam"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center justify-center gap-2 px-6 py-4 rounded-xl bg-white hover:bg-[#F8FAFC] text-[#1E2B58] border border-[#E2E8F0] font-bold text-sm sm:text-base shadow-sm hover:border-[#00B894] transition-all"
              >
                <Send className="w-4 h-4 text-[#00B894]" />
                <span>Join Telegram</span>
              </a>
            </div>

            {/* Trust points */}
            <div className="flex flex-wrap items-center gap-6 pt-3 text-xs font-semibold text-[#64748B]">
              <div className="flex items-center gap-1.5">
                <CheckCircle2 className="w-4 h-4 text-[#00B894]" />
                <span>Free Starter Access</span>
              </div>
              <div className="flex items-center gap-1.5">
                <CheckCircle2 className="w-4 h-4 text-[#00B894]" />
                <span>Natural & Social Streams</span>
              </div>
              <div className="flex items-center gap-1.5">
                <CheckCircle2 className="w-4 h-4 text-[#00B894]" />
                <span>Telebirr & CBE Ready</span>
              </div>
            </div>

          </div>

          {/* Right Column: Dynamic Vector Illustration (Stylistically matched to video frame 002) */}
          <div className="lg:col-span-6 relative flex items-center justify-center">
            
            <div className="relative w-full max-w-[520px] aspect-4/3 flex items-center justify-center">
              
              {/* Soft Background Blob */}
              <div className="absolute inset-0 bg-gradient-to-tr from-[#E8F8F5] via-white to-[#F1F5F9] rounded-[40px] shadow-dribbble -rotate-1 border border-[#E2E8F0]/60" />

              {/* Vector Composition */}
              <svg viewBox="0 0 520 420" className="w-full h-auto drop-shadow-md relative z-10">
                <defs>
                  <linearGradient id="bulbGrad" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor="#00B894" />
                    <stop offset="100%" stopColor="#00A381" />
                  </linearGradient>
                  <linearGradient id="rocketGrad" x1="0" y1="0" x2="1" y2="1">
                    <stop offset="0%" stopColor="#1E2B58" />
                    <stop offset="100%" stopColor="#2A3B72" />
                  </linearGradient>
                  <linearGradient id="tableGrad" x1="0" y1="0" x2="1" y2="0">
                    <stop offset="0%" stopColor="#E2E8F0" />
                    <stop offset="100%" stopColor="#CBD5E1" />
                  </linearGradient>
                </defs>

                {/* Desk Surface */}
                <ellipse cx="260" cy="385" rx="200" ry="12" fill="#E2E8F0" />
                <rect x="120" y="270" width="280" height="10" rx="5" fill="url(#tableGrad)" />
                <line x1="160" y1="280" x2="145" y2="380" stroke="#94A3B8" strokeWidth="6" strokeLinecap="round" />
                <line x1="360" y1="280" x2="375" y2="380" stroke="#94A3B8" strokeWidth="6" strokeLinecap="round" />

                {/* Collaborative Students Group (Matching the team illustration in Dribbble Video) */}
                {/* Student 1 (Left - Green Shirt with Laptop) */}
                <circle cx="160" cy="210" r="16" fill="#F8CBA6" />
                <path d="M148 205 C148 190, 172 190, 172 205 Z" fill="#1E2B58" />
                <path d="M140 255 C140 230, 180 230, 180 255 L180 320 L140 320 Z" fill="#00B894" />
                <path d="M165 245 L200 270 L195 275 L160 250 Z" fill="#00A381" />

                {/* Laptop on desk */}
                <rect x="180" y="245" width="40" height="25" rx="3" fill="#1E2B58" />
                <rect x="175" y="268" width="50" height="4" rx="2" fill="#CBD5E1" />
                <circle cx="200" cy="257" r="4" fill="#00B894" />

                {/* Student 2 (Center - Pointing up to Idea Bulb) */}
                <circle cx="260" cy="180" r="18" fill="#F8CBA6" />
                <path d="M246 175 C246 160, 274 160, 274 175 Z" fill="#1E2B58" />
                <path d="M240 225 C240 205, 280 205, 280 225 L285 300 L235 300 Z" fill="#1E2B58" />
                {/* Pointing Arm pointing toward puzzle bulb */}
                <line x1="275" y1="215" x2="310" y2="155" stroke="#F8CBA6" strokeWidth="8" strokeLinecap="round" />

                {/* Giant Glowing Puzzle Idea Bulb (Centerpiece of the Dribbble video) */}
                <g className="animate-pulse" style={{ animationDuration: "3s" }}>
                  <circle cx="330" cy="120" r="38" fill="url(#bulbGrad)" filter="drop-shadow(0 8px 16px rgba(0,184,148,0.3))" />
                  <rect x="320" y="155" width="20" height="12" rx="3" fill="#1E2B58" />
                  <line x1="324" y1="162" x2="336" y2="162" stroke="#CBD5E1" strokeWidth="2" />
                  <line x1="324" y1="165" x2="336" y2="165" stroke="#CBD5E1" strokeWidth="2" />
                  {/* Puzzle outline pattern inside bulb */}
                  <path d="M320 105 Q330 95, 340 105 Q345 115, 335 125 Q325 135, 330 140" stroke="#FFFFFF" strokeWidth="2.5" strokeLinecap="round" fill="none" opacity="0.8" />
                  <circle cx="330" cy="120" r="8" fill="#FFFFFF" opacity="0.9" />
                </g>

                {/* Student 3 (Right - Mentor / Educator Standing) */}
                <circle cx="370" cy="195" r="17" fill="#F8CBA6" />
                <path d="M356 190 C356 175, 384 175, 384 190 Z" fill="#1E2B58" />
                <path d="M350 240 C350 215, 390 215, 390 240 L395 340 L345 340 Z" fill="#2A3B72" />
                <line x1="355" y1="230" x2="330" y2="260" stroke="#F8CBA6" strokeWidth="7" strokeLinecap="round" />

                {/* Floating Rocket Launch Graphic (from Dribbble video) */}
                <g className="animate-bounce" style={{ animationDuration: "4s" }}>
                  <path d="M410 90 L435 65 Q450 75, 445 95 L420 120 Z" fill="url(#rocketGrad)" />
                  <circle cx="430" cy="85" r="5" fill="#00B894" />
                  <path d="M422 118 L415 130 L425 125 Z" fill="#EF4444" />
                  <path d="M428 114 L432 128 L438 122 Z" fill="#F59E0B" />
                  <circle cx="425" cy="135" r="3" fill="#E2E8F0" opacity="0.6" />
                  <circle cx="430" cy="145" r="4" fill="#E2E8F0" opacity="0.4" />
                </g>

                {/* Floating Analytics Progress Card (Left Top) */}
                <g>
                  <rect x="60" y="110" width="95" height="55" rx="10" fill="#FFFFFF" filter="drop-shadow(0 8px 16px rgba(30,43,88,0.08))" />
                  <rect x="72" y="122" width="40" height="5" rx="2.5" fill="#1E2B58" />
                  <rect x="72" y="132" width="60" height="4" rx="2" fill="#94A3B8" />
                  <rect x="72" y="148" width="10" height="10" rx="2" fill="#E8F8F5" />
                  <text x="88" y="156" fill="#00B894" fontSize="9" fontWeight="bold">94% Pass</text>
                </g>

                {/* Floating Score Rank Pill (Right Bottom) */}
                <g>
                  <rect x="360" y="270" width="130" height="42" rx="12" fill="#FFFFFF" filter="drop-shadow(0 8px 20px rgba(30,43,88,0.1))" />
                  <circle cx="380" cy="291" r="10" fill="#E8F8F5" />
                  <text x="376" y="295" fill="#00B894" fontSize="12" fontWeight="bold">★</text>
                  <text x="398" y="287" fill="#1E2B58" fontSize="10" fontWeight="bold">642 / 700 Score</text>
                  <text x="398" y="298" fill="#64748B" fontSize="8">AAU Medicine Cutoff</text>
                </g>
              </svg>

            </div>

          </div>

        </div>

      </div>
    </section>
  );
}
