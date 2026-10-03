"use client";

import React from "react";
import Image from "next/image";
import { 
  WifiOff, 
  DownloadCloud, 
  HardDrive, 
  BatteryCharging, 
  CheckCircle, 
  Smartphone,
  ShieldAlert,
  Sparkles,
  Zap,
  ArrowRight
} from "lucide-react";

export function OfflineShowcase() {
  const steps = [
    {
      num: "01",
      icon: DownloadCloud,
      title: "One-Click Unit Download",
      desc: "Download high-definition video lessons and PDF unit notes when connected to school, cafe, or home Wi-Fi.",
      highlight: "Saves mobile data",
    },
    {
      num: "02",
      icon: HardDrive,
      title: "Encrypted Device Storage",
      desc: "Lessons are compressed and securely saved to your phone's storage, taking minimal device space (~18MB per chapter).",
      highlight: "Lightweight & Fast",
    },
    {
      num: "03",
      icon: WifiOff,
      title: "Study Anywhere in Ethiopia",
      desc: "Turn off mobile data completely and take timed mock exams, rewatch difficult lessons, and review notes with zero internet connection.",
      highlight: "100% Data-Free",
    },
  ];

  return (
    <section id="offline" className="py-24 bg-[#08080b] border-t border-[#1a1a22] relative overflow-hidden">
      {/* Background glow */}
      <div className="absolute top-1/2 right-10 -translate-y-1/2 w-[500px] h-[500px] bg-[#10b981]/10 rounded-full blur-[160px] pointer-events-none" />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 items-center">
          
          {/* Left: Explainer Copy & Steps */}
          <div className="lg:col-span-6 space-y-6">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-[#10b981]/15 text-[#10b981] border border-[#10b981]/30 text-xs font-bold">
              <WifiOff className="w-3.5 h-3.5" />
              <span>Offline Study Architecture</span>
            </div>

            <h2 className="text-3xl sm:text-4xl font-extrabold text-white tracking-tight">
              Study Without Limits — Even With No Internet Connection
            </h2>

            <p className="text-sm sm:text-base text-[#94a3b8] leading-relaxed">
              We know internet across Ethiopian regions can be slow, expensive, or interrupted. Mirkuz was architected so Grade 12 students never miss a single day of study.
            </p>

            <div className="space-y-4 pt-2">
              {steps.map((step, idx) => {
                const Icon = step.icon;
                return (
                  <div
                    key={idx}
                    className="p-5 rounded-2xl bg-[#111116] border border-[#1e1e26] hover:border-[#2a2a35] transition-all flex items-start gap-4"
                  >
                    <div className="w-10 h-10 rounded-xl bg-[#10b981]/20 border border-[#10b981]/40 flex items-center justify-center shrink-0 text-[#10b981] font-bold text-xs">
                      {step.num}
                    </div>
                    <div className="space-y-1">
                      <div className="flex items-center gap-2">
                        <h3 className="text-base font-bold text-white">{step.title}</h3>
                        <span className="text-[10px] px-2 py-0.5 rounded-full bg-[#10b981]/10 text-[#10b981] font-semibold border border-[#10b981]/20">
                          {step.highlight}
                        </span>
                      </div>
                      <p className="text-xs text-[#94a3b8] leading-relaxed">
                        {step.desc}
                      </p>
                    </div>
                  </div>
                );
              })}
            </div>

            {/* Savings Callout */}
            <div className="p-4 rounded-2xl bg-gradient-to-r from-[#10b981]/15 via-[#059669]/15 to-transparent border border-[#10b981]/30 flex items-center gap-3.5">
              <Sparkles className="w-5 h-5 text-[#10b981] shrink-0" />
              <p className="text-xs text-[#e2e8f0]">
                <strong className="text-white font-semibold">Average Savings:</strong> Students save over <span className="text-[#10b981] font-bold">15 GB</span> of mobile internet packages every month.
              </p>
            </div>
          </div>

          {/* Right: Phone Frame with Real Offline Screen Screenshot */}
          <div className="lg:col-span-6 flex justify-center">
            <div className="relative">
              
              {/* Radial glow */}
              <div className="absolute inset-0 bg-gradient-to-tr from-[#10b981]/25 to-transparent rounded-[50px] blur-2xl pointer-events-none" />

              {/* Smartphone Container */}
              <div className="relative w-[280px] sm:w-[320px] rounded-[46px] p-3 bg-gradient-to-b from-[#3a3a46] via-[#1c1c24] to-[#111116] shadow-2xl shadow-black/90 border border-[#484858]">
                
                <div className="relative rounded-[36px] bg-[#050505] overflow-hidden border border-[#181820] shadow-inner aspect-[9/18.5]">
                  <Image
                    src="/screenshots/mirkuz_offline.jpg"
                    alt="Mirkuz Offline Downloads Manager"
                    fill
                    sizes="(max-width: 640px) 280px, 320px"
                    className="object-cover"
                  />
                  <div className="absolute inset-0 bg-gradient-to-tr from-transparent via-white/[0.03] to-transparent pointer-events-none" />
                </div>

                <div className="w-24 h-1 bg-white/30 rounded-full mx-auto mt-2" />
              </div>

              {/* Floating Badge */}
              <div className="absolute -bottom-5 -right-4 sm:-right-6 bg-[#14141a]/95 backdrop-blur-md p-3.5 rounded-2xl border border-[#2c2c3a] shadow-2xl flex items-center gap-3">
                <div className="w-9 h-9 rounded-xl bg-[#10b981]/20 flex items-center justify-center text-[#10b981]">
                  <WifiOff className="w-4 h-4" />
                </div>
                <div>
                  <div className="text-xs font-bold text-white">0 MB Internet Used</div>
                  <div className="text-[10px] text-[#94a3b8]">100% Offline Exam Mode</div>
                </div>
              </div>

            </div>
          </div>

        </div>

      </div>
    </section>
  );
}
