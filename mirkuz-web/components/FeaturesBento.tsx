"use client";

import React from "react";
import { 
  WifiOff, 
  Play, 
  Clock, 
  Bot, 
  BarChart3, 
  FileCheck, 
  Download, 
  Zap, 
  ShieldCheck, 
  Sparkles,
  Smartphone,
  Cpu
} from "lucide-react";

export function FeaturesBento() {
  return (
    <section id="features" className="py-24 bg-[#09090d] relative overflow-hidden">
      {/* Glow */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[500px] bg-[#10b981]/10 rounded-full blur-[160px] pointer-events-none" />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        
        {/* Section Header */}
        <div className="text-center max-w-3xl mx-auto mb-16 space-y-4">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-[#38bdf8]/10 text-[#38bdf8] border border-[#38bdf8]/25 text-xs font-semibold">
            <Zap className="w-3.5 h-3.5" />
            <span>Built Specifically for Ethiopian Students</span>
          </div>
          <h2 className="text-3xl sm:text-4xl font-extrabold text-white tracking-tight">
            Engineered for High Performance, Low Bandwidth & Maximum Exam Scores
          </h2>
          <p className="text-sm sm:text-base text-[#94a3b8]">
            Studying in Ethiopia comes with unique connectivity challenges. Mirkuz was engineered from the ground up to guarantee uninterrupted learning on any smartphone.
          </p>
        </div>

        {/* Bento Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-12 gap-6">
          
          {/* Card 1: Offline First (Col span 7) */}
          <div className="lg:col-span-7 rounded-3xl bg-[#111115] border border-[#1a1a20] p-8 relative overflow-hidden flex flex-col justify-between hover:border-[#2c2c35] transition-all group">
            <div className="relative z-10 space-y-4">
              <div className="w-12 h-12 rounded-2xl bg-[#10b981]/20 border border-[#10b981]/30 flex items-center justify-center text-[#10b981]">
                <WifiOff className="w-6 h-6" />
              </div>
              <div>
                <span className="text-xs font-bold uppercase tracking-wider text-[#10b981]">Zero Data Study</span>
                <h3 className="text-2xl font-bold text-white mt-1">
                  100% Encrypted Offline Learning Mode
                </h3>
              </div>
              <p className="text-sm text-[#94a3b8] leading-relaxed max-w-lg">
                Download entire subject units, HD video explanations, and PDF study sheets directly onto your phone’s internal storage. Turn off mobile data and study anytime without spending a single Birr on internet packages.
              </p>
            </div>

            {/* Offline Visual Interactive preview */}
            <div className="mt-8 pt-6 border-t border-[#1a1a20] grid grid-cols-3 gap-3">
              <div className="p-3 rounded-xl bg-[#18181e] border border-[#2c2c35] text-center">
                <span className="text-[10px] text-[#94a3b8] block">Step 1</span>
                <span className="text-xs font-bold text-white">Download on Wi-Fi</span>
              </div>
              <div className="p-3 rounded-xl bg-[#18181e] border border-[#2c2c35] text-center">
                <span className="text-[10px] text-[#94a3b8] block">Step 2</span>
                <span className="text-xs font-bold text-[#10b981]">Toggle Offline</span>
              </div>
              <div className="p-3 rounded-xl bg-[#18181e] border border-[#2c2c35] text-center">
                <span className="text-[10px] text-[#94a3b8] block">Step 3</span>
                <span className="text-xs font-bold text-white">0MB Data Used</span>
              </div>
            </div>
          </div>

          {/* Card 2: Adaptive HLS Streaming (Col span 5) */}
          <div className="lg:col-span-5 rounded-3xl bg-[#111115] border border-[#1a1a20] p-8 relative overflow-hidden flex flex-col justify-between hover:border-[#2c2c35] transition-all group">
            <div className="relative z-10 space-y-4">
              <div className="w-12 h-12 rounded-2xl bg-[#38bdf8]/20 border border-[#38bdf8]/30 flex items-center justify-center text-[#38bdf8]">
                <Play className="w-6 h-6 fill-current" />
              </div>
              <div>
                <span className="text-xs font-bold uppercase tracking-wider text-[#38bdf8]">Adaptive Streaming</span>
                <h3 className="text-2xl font-bold text-white mt-1">
                  Low-Latency HLS Video Engine
                </h3>
              </div>
              <p className="text-sm text-[#94a3b8] leading-relaxed">
                Smart adaptive bitrate automatically scales video resolution between 360p and 1080p, guaranteeing zero buffering even on congested 3G Ethiopian mobile connections.
              </p>
            </div>

            <div className="mt-6 flex items-center gap-2">
              <span className="px-2.5 py-1 rounded-md bg-[#18181e] text-[11px] font-mono text-[#38bdf8] border border-[#2c2c35]">360p Eco</span>
              <span className="px-2.5 py-1 rounded-md bg-[#18181e] text-[11px] font-mono text-[#38bdf8] border border-[#2c2c35]">720p HD</span>
              <span className="px-2.5 py-1 rounded-md bg-[#18181e] text-[11px] font-mono text-[#38bdf8] border border-[#2c2c35]">1.5x / 2.0x Speed</span>
            </div>
          </div>

          {/* Card 3: Timed Mock Exams (Col span 4) */}
          <div className="lg:col-span-4 rounded-3xl bg-[#111115] border border-[#1a1a20] p-7 flex flex-col justify-between hover:border-[#2c2c35] transition-all">
            <div className="space-y-3">
              <div className="w-11 h-11 rounded-xl bg-[#f59e0b]/20 border border-[#f59e0b]/30 flex items-center justify-center text-[#f59e0b]">
                <Clock className="w-5 h-5" />
              </div>
              <h3 className="text-lg font-bold text-white">
                Simulated National Entrance Exams
              </h3>
              <p className="text-xs text-[#94a3b8] leading-relaxed">
                Experience the authentic exam atmosphere with real-time countdown timers, question review flags, automatic scoring, and national percentile ranks.
              </p>
            </div>
            <div className="mt-4 pt-3 border-t border-[#1a1a20] text-xs font-semibold text-[#f59e0b]">
              ✓ Instant step-by-step video solutions
            </div>
          </div>

          {/* Card 4: AI Concept Tutor Mascot (Col span 4) */}
          <div className="lg:col-span-4 rounded-3xl bg-[#111115] border border-[#1a1a20] p-7 flex flex-col justify-between hover:border-[#2c2c35] transition-all">
            <div className="space-y-3">
              <div className="w-11 h-11 rounded-xl bg-[#8b5cf6]/20 border border-[#8b5cf6]/30 flex items-center justify-center text-[#8b5cf6]">
                <Bot className="w-5 h-5" />
              </div>
              <h3 className="text-lg font-bold text-white">
                Mirkuz AI Concept Tutor
              </h3>
              <p className="text-xs text-[#94a3b8] leading-relaxed">
                Stuck on a difficult calculus derivation or physics circuit? Tap the Mirkuz Mascot AI assistant for instant hints and underlying conceptual explanations.
              </p>
            </div>
            <div className="mt-4 pt-3 border-t border-[#1a1a20] text-xs font-semibold text-[#8b5cf6]">
              ✓ 24/7 personalized study assistance
            </div>
          </div>

          {/* Card 5: Weak-Area Diagnostic (Col span 4) */}
          <div className="lg:col-span-4 rounded-3xl bg-[#111115] border border-[#1a1a20] p-7 flex flex-col justify-between hover:border-[#2c2c35] transition-all">
            <div className="space-y-3">
              <div className="w-11 h-11 rounded-xl bg-[#10b981]/20 border border-[#10b981]/30 flex items-center justify-center text-[#10b981]">
                <BarChart3 className="w-5 h-5" />
              </div>
              <h3 className="text-lg font-bold text-white">
                Weak-Area Diagnostic Analytics
              </h3>
              <p className="text-xs text-[#94a3b8] leading-relaxed">
                See exactly which units and topics pull your scores down. Track mastery percentages chapter-by-chapter and focus your revision where it matters most.
              </p>
            </div>
            <div className="mt-4 pt-3 border-t border-[#1a1a20] text-xs font-semibold text-[#10b981]">
              ✓ Personalized score improvement roadmap
            </div>
          </div>

        </div>

      </div>
    </section>
  );
}
