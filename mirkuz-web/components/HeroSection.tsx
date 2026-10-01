"use client";

import React from "react";
import { 
  Download, 
  Send, 
  Play, 
  CheckCircle2, 
  Star, 
  Sparkles, 
  Clock, 
  FileText, 
  WifiOff, 
  ShieldCheck, 
  ArrowRight,
  TrendingUp,
  BrainCircuit,
  BookOpen
} from "lucide-react";

export function HeroSection() {
  return (
    <section className="relative pt-32 pb-20 md:pt-40 md:pb-28 overflow-hidden bg-grid-pattern">
      {/* Radial Glow Backdrops */}
      <div className="absolute top-1/4 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[600px] h-[600px] bg-[#10b981]/15 rounded-full blur-[140px] pointer-events-none" />
      <div className="absolute top-1/3 right-10 w-[400px] h-[400px] bg-[#f59e0b]/10 rounded-full blur-[120px] pointer-events-none" />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-8 items-center">
          
          {/* Left Column: Copy & CTAs */}
          <div className="lg:col-span-7 flex flex-col items-start text-left space-y-6">
            
            {/* National Exam Badge */}
            <div className="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-[#18181e] border border-[#2c2c35] text-xs font-medium text-[#94a3b8] shadow-inner">
              <span className="flex h-2 w-2 relative">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-[#10b981] opacity-75"></span>
                <span className="relative inline-flex rounded-full h-2 w-2 bg-[#10b981]"></span>
              </span>
              <span>🇪🇹 Aligned with Ethiopian Grade 12 National Curriculum</span>
            </div>

            {/* Headline */}
            <h1 className="text-4xl sm:text-5xl lg:text-6xl font-extrabold tracking-tight text-white leading-[1.15]">
              Master Your Grade 12{" "}
              <span className="bg-gradient-to-r from-[#10b981] via-[#34d399] to-[#f59e0b] bg-clip-text text-transparent">
                Entrance Exam
              </span>{" "}
              with Mirkuz
            </h1>

            {/* Sub-headline */}
            <p className="text-base sm:text-lg text-[#94a3b8] max-w-2xl leading-relaxed">
              Ethiopia’s premier exam prep platform. HD video lessons by top Ethiopian educators, downloadable PDF chapter summaries, realistic timed mock exams with instant step-by-step solutions, and <strong className="text-white font-semibold">100% data-free offline study</strong>.
            </p>

            {/* Primary Action Buttons */}
            <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-4 w-full sm:w-auto pt-2">
              <a
                href="https://play.google.com/store/apps"
                target="_blank"
                rel="noopener noreferrer"
                className="group relative flex items-center justify-center gap-3 px-6 py-4 rounded-xl bg-gradient-to-r from-[#10b981] to-[#059669] text-white font-semibold shadow-xl shadow-[#10b981]/25 hover:shadow-2xl hover:shadow-[#10b981]/40 hover:scale-[1.02] active:scale-[0.98] transition-all"
              >
                <div className="w-6 h-6 flex items-center justify-center">
                  <svg className="w-5 h-5 fill-current" viewBox="0 0 24 24">
                    <path d="M3.609 1.814L13.792 12 3.61 22.186a2.41 2.41 0 0 1-.61-.926V2.74c.15-.36.368-.68.61-.926zm11.597 11.598l2.584 2.585-12.01 6.864 9.426-9.449zm2.584-2.824l-2.584 2.584L5.78 3.723l12.01 6.865zm1.536.878l2.94 1.68a1.2 1.2 0 0 1 0 2.086l-2.94 1.68-2.197-2.197 2.197-2.249z" />
                  </svg>
                </div>
                <div className="flex flex-col text-left">
                  <span className="text-[10px] uppercase tracking-wider text-emerald-100 font-bold">Available on</span>
                  <span className="text-base font-bold leading-tight">Google Play Store</span>
                </div>
                <ArrowRight className="w-4 h-4 ml-1 opacity-70 group-hover:translate-x-1 transition-transform" />
              </a>

              <a
                href="https://t.me/mirkuz_exam"
                target="_blank"
                rel="noopener noreferrer"
                className="flex items-center justify-center gap-3 px-6 py-4 rounded-xl bg-[#18181e] hover:bg-[#25252e] text-white font-semibold border border-[#2c2c35] hover:border-[#38bdf8]/50 shadow-lg hover:shadow-xl hover:shadow-[#38bdf8]/10 transition-all active:scale-[0.98]"
              >
                <Send className="w-5 h-5 text-[#38bdf8]" />
                <div className="flex flex-col text-left">
                  <span className="text-[10px] uppercase tracking-wider text-[#94a3b8] font-bold">Join Community</span>
                  <span className="text-base font-bold leading-tight text-white">Telegram Channel</span>
                </div>
              </a>
            </div>

            {/* Quick interactive test jump & ratings */}
            <div className="flex flex-wrap items-center gap-6 pt-2 text-xs text-[#94a3b8]">
              <div className="flex items-center gap-1.5 text-[#f59e0b]">
                <div className="flex">
                  {[...Array(5)].map((_, i) => (
                    <Star key={i} className="w-4 h-4 fill-current text-[#f59e0b]" />
                  ))}
                </div>
                <span className="font-bold text-white ml-1">4.9 / 5</span>
                <span className="text-[#71717a]">(3,400+ Grade 12 Reviews)</span>
              </div>

              <a
                href="#exam-demo"
                className="inline-flex items-center gap-1.5 text-[#10b981] font-semibold hover:underline"
              >
                <Play className="w-3.5 h-3.5 fill-current" />
                <span>Try Live Entrance Question Demo</span>
              </a>
            </div>

            {/* Value Pillars Pills */}
            <div className="grid grid-cols-2 sm:grid-cols-3 gap-3 w-full pt-4 border-t border-[#1a1a20]">
              <div className="flex items-center gap-2.5 p-2.5 rounded-lg bg-[#111115]/80 border border-[#1a1a20]">
                <WifiOff className="w-4 h-4 text-[#10b981]" />
                <span className="text-xs font-medium text-[#d4d4d8]">100% Offline Study</span>
              </div>
              <div className="flex items-center gap-2.5 p-2.5 rounded-lg bg-[#111115]/80 border border-[#1a1a20]">
                <Clock className="w-4 h-4 text-[#f59e0b]" />
                <span className="text-xs font-medium text-[#d4d4d8]">Timed Mock Exams</span>
              </div>
              <div className="flex items-center gap-2.5 p-2.5 rounded-lg bg-[#111115]/80 border border-[#1a1a20] col-span-2 sm:col-span-1">
                <ShieldCheck className="w-4 h-4 text-[#38bdf8]" />
                <span className="text-xs font-medium text-[#d4d4d8]">Telebirr & Chapa Pay</span>
              </div>
            </div>

          </div>

          {/* Right Column: High-Fidelity Interactive App Mockup */}
          <div className="lg:col-span-5 relative flex justify-center lg:justify-end">
            {/* Glowing Aura Behind Phone */}
            <div className="absolute inset-0 bg-gradient-to-tr from-[#10b981]/20 via-[#f59e0b]/15 to-transparent rounded-3xl blur-2xl transform rotate-3 scale-95 pointer-events-none" />

            {/* Smartphone Frame Container */}
            <div className="relative w-full max-w-[340px] sm:max-w-[370px] rounded-[42px] p-3.5 bg-gradient-to-b from-[#2c2c35] via-[#1a1a20] to-[#111115] shadow-2xl shadow-black/80 border border-[#3f3f4e]">
              
              {/* Inner Bezel Screen */}
              <div className="relative rounded-[32px] bg-[#050505] overflow-hidden border border-[#18181e] text-white">
                
                {/* Phone Status Bar */}
                <div className="px-5 pt-3 pb-2 flex items-center justify-between text-[11px] text-[#94a3b8] bg-[#09090b]">
                  <span className="font-semibold text-white">09:41</span>
                  {/* Speaker Notch */}
                  <div className="w-20 h-4 bg-[#18181e] rounded-full mx-auto border border-[#2c2c35]/50 flex items-center justify-center">
                    <div className="w-3 h-3 rounded-full bg-[#10b981]/40"></div>
                  </div>
                  <div className="flex items-center gap-1.5">
                    <span className="text-[10px] font-bold text-[#10b981]">4G</span>
                    <div className="w-4 h-2.5 border border-[#94a3b8] rounded-sm p-0.5 flex items-center">
                      <div className="w-full h-full bg-[#10b981] rounded-2xs"></div>
                    </div>
                  </div>
                </div>

                {/* App Header */}
                <div className="px-4 py-3 bg-[#111115] border-b border-[#1a1a20] flex items-center justify-between">
                  <div className="flex items-center gap-2">
                    <div className="w-7 h-7 rounded-lg bg-[#10b981] flex items-center justify-center text-white text-xs font-bold">
                      M
                    </div>
                    <div>
                      <h4 className="text-xs font-bold text-white">Mathematics (Grade 12)</h4>
                      <p className="text-[10px] text-[#94a3b8]">Unit 3: Integral Calculus</p>
                    </div>
                  </div>
                  <span className="px-2 py-0.5 rounded-full text-[10px] font-semibold bg-[#10b981]/20 text-[#10b981] border border-[#10b981]/40">
                    Offline Saved
                  </span>
                </div>

                {/* Simulated Video Player / Lesson Card */}
                <div className="relative aspect-video bg-[#18181e] flex flex-col justify-between p-3 overflow-hidden group">
                  <div className="absolute inset-0 bg-gradient-to-t from-black/80 via-black/20 to-black/60" />
                  
                  {/* Top video controls */}
                  <div className="relative z-10 flex items-center justify-between text-[10px] text-white/90">
                    <span className="px-2 py-0.5 rounded bg-black/60 backdrop-blur font-mono">
                      Lesson 4 of 12
                    </span>
                    <span className="px-2 py-0.5 rounded bg-[#10b981]/80 font-bold text-black text-[9px]">
                      HLS HD 720p
                    </span>
                  </div>

                  {/* Center Play Indicator */}
                  <div className="relative z-10 self-center w-12 h-12 rounded-full bg-[#10b981]/90 backdrop-blur flex items-center justify-center text-white shadow-lg shadow-[#10b981]/40 cursor-pointer hover:scale-110 transition-transform">
                    <Play className="w-5 h-5 fill-current ml-0.5" />
                  </div>

                  {/* Bottom Video Progress */}
                  <div className="relative z-10 space-y-1">
                    <div className="flex justify-between text-[10px] text-white/80 font-mono">
                      <span>14:20</span>
                      <span>28:45</span>
                    </div>
                    <div className="w-full h-1.5 bg-white/20 rounded-full overflow-hidden">
                      <div className="w-[52%] h-full bg-[#10b981] rounded-full"></div>
                    </div>
                  </div>
                </div>

                {/* Live Mock Quiz Card Inside App Mockup */}
                <div className="p-4 space-y-3 bg-[#0c0c10]">
                  
                  {/* Timer & Question Number */}
                  <div className="flex items-center justify-between text-xs">
                    <div className="flex items-center gap-1.5 text-[#f59e0b] font-bold">
                      <Clock className="w-3.5 h-3.5" />
                      <span>Time Left: 01:45</span>
                    </div>
                    <span className="text-[11px] text-[#94a3b8] font-medium">Question 14/50</span>
                  </div>

                  {/* Question Box */}
                  <div className="p-3 rounded-xl bg-[#18181e] border border-[#2c2c35]">
                    <p className="text-xs font-semibold text-white leading-relaxed">
                      Evaluate the integral: <code className="text-[#38bdf8] font-mono">∫ (3x² + 4x - 5) dx</code>
                    </p>
                  </div>

                  {/* Answer Options */}
                  <div className="space-y-1.5 text-xs">
                    <div className="p-2.5 rounded-lg bg-[#10b981]/20 border border-[#10b981] text-[#10b981] font-medium flex items-center justify-between">
                      <span>A. x³ + 2x² - 5x + C</span>
                      <CheckCircle2 className="w-4 h-4 text-[#10b981]" />
                    </div>
                    <div className="p-2.5 rounded-lg bg-[#111115] border border-[#1a1a20] text-[#94a3b8]">
                      <span>B. 6x + 4 + C</span>
                    </div>
                  </div>

                  {/* Mirkuz Mascot AI Tutor Pill */}
                  <div className="p-2.5 rounded-xl bg-gradient-to-r from-[#18181e] to-[#111115] border border-[#2c2c35] flex items-center gap-3 shadow-md">
                    <div className="w-8 h-8 rounded-full bg-[#10b981]/20 border border-[#10b981]/50 flex items-center justify-center shrink-0">
                      <BrainCircuit className="w-4 h-4 text-[#10b981]" />
                    </div>
                    <div className="text-[11px]">
                      <span className="font-bold text-white block">Mirkuz AI Step-by-Step</span>
                      <span className="text-[#94a3b8] text-[10px]">Power rule applied: 3(x³/3) + 4(x²/2) - 5x + C</span>
                    </div>
                  </div>

                </div>

                {/* App Bottom Navigation */}
                <div className="px-6 py-2.5 bg-[#09090b] border-t border-[#1a1a20] flex items-center justify-around text-[#71717a]">
                  <div className="flex flex-col items-center gap-1 text-[#10b981]">
                    <BookOpen className="w-4 h-4" />
                    <span className="text-[9px] font-semibold">Courses</span>
                  </div>
                  <div className="flex flex-col items-center gap-1">
                    <Clock className="w-4 h-4" />
                    <span className="text-[9px]">Mock Exam</span>
                  </div>
                  <div className="flex flex-col items-center gap-1">
                    <TrendingUp className="w-4 h-4" />
                    <span className="text-[9px]">Analytics</span>
                  </div>
                  <div className="flex flex-col items-center gap-1">
                    <WifiOff className="w-4 h-4" />
                    <span className="text-[9px]">Saved</span>
                  </div>
                </div>

              </div>
            </div>

            {/* Floating Live Indicator Badge */}
            <div className="hidden sm:flex absolute -bottom-6 -left-6 bg-[#18181e]/95 backdrop-blur-md p-3.5 rounded-2xl border border-[#2c2c35] shadow-2xl items-center gap-3">
              <div className="w-10 h-10 rounded-xl bg-[#10b981]/20 flex items-center justify-center text-[#10b981]">
                <TrendingUp className="w-5 h-5" />
              </div>
              <div>
                <div className="text-xs font-bold text-white">+18.5% Score Boost</div>
                <div className="text-[11px] text-[#94a3b8]">Average National Exam Improvement</div>
              </div>
            </div>

          </div>

        </div>
      </div>
    </section>
  );
}
