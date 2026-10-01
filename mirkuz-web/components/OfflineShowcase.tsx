"use client";

import React from "react";
import { 
  WifiOff, 
  DownloadCloud, 
  HardDrive, 
  BatteryCharging, 
  CheckCircle, 
  Smartphone,
  ShieldAlert,
  Sparkles
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
      title: "Encrypted Hive Storage",
      desc: "Lessons are compressed and securely saved to your phone's storage, taking minimal device space (~18MB per chapter).",
      highlight: "Lightweight & Fast",
    },
    {
      num: "03",
      icon: WifiOff,
      title: "Study Anywhere in Ethiopia",
      desc: "Turn off mobile data and take timed mock exams, rewatch difficult lessons, and review notes with zero internet connection.",
      highlight: "100% Data-Free",
    },
  ];

  return (
    <section id="offline" className="py-24 bg-[#09090d] border-t border-[#1a1a20] relative">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
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
              We know internet in Ethiopia can be slow, expensive, or interrupted. Mirkuz was architected so Grade 12 students never miss a day of study.
            </p>

            <div className="space-y-4 pt-4">
              {steps.map((step, idx) => {
                const Icon = step.icon;
                return (
                  <div
                    key={idx}
                    className="p-5 rounded-2xl bg-[#111115] border border-[#1a1a20] hover:border-[#2c2c35] transition-all flex items-start gap-4"
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
          </div>

          {/* Right: Visual Storage & Download Card */}
          <div className="lg:col-span-6 relative">
            <div className="rounded-3xl bg-[#111115] border border-[#2c2c35] p-8 shadow-2xl relative overflow-hidden space-y-6">
              
              <div className="flex items-center justify-between border-b border-[#1a1a20] pb-4">
                <div className="flex items-center gap-2.5">
                  <div className="w-3 h-3 rounded-full bg-[#10b981]"></div>
                  <span className="text-sm font-bold text-white">Mirkuz Offline Storage Manager</span>
                </div>
                <span className="text-xs font-mono text-[#10b981]">Status: 100% Ready</span>
              </div>

              {/* Downloaded Subject Items */}
              <div className="space-y-3">
                
                {/* Item 1 */}
                <div className="p-4 rounded-xl bg-[#18181e] border border-[#2c2c35] flex items-center justify-between">
                  <div className="space-y-1">
                    <span className="text-xs font-bold text-white block">Physics — Unit 3: Electromagnetism</span>
                    <span className="text-[10px] text-[#94a3b8]">6 Lessons + PDF Summary • 84 MB</span>
                  </div>
                  <span className="px-2.5 py-1 rounded-lg bg-[#10b981]/20 text-[#10b981] text-xs font-bold border border-[#10b981]/40 flex items-center gap-1">
                    <CheckCircle className="w-3.5 h-3.5" />
                    <span>Saved</span>
                  </span>
                </div>

                {/* Item 2 */}
                <div className="p-4 rounded-xl bg-[#18181e] border border-[#2c2c35] flex items-center justify-between">
                  <div className="space-y-1">
                    <span className="text-xs font-bold text-white block">Chemistry — Unit 4: Electrochemistry</span>
                    <span className="text-[10px] text-[#94a3b8]">8 Lessons + Past Exams • 96 MB</span>
                  </div>
                  <span className="px-2.5 py-1 rounded-lg bg-[#10b981]/20 text-[#10b981] text-xs font-bold border border-[#10b981]/40 flex items-center gap-1">
                    <CheckCircle className="w-3.5 h-3.5" />
                    <span>Saved</span>
                  </span>
                </div>

                {/* Item 3 */}
                <div className="p-4 rounded-xl bg-[#18181e] border border-[#2c2c35] flex items-center justify-between">
                  <div className="space-y-1">
                    <span className="text-xs font-bold text-white block">Mathematics — 2016 E.C. Full Mock Exam</span>
                    <span className="text-[10px] text-[#94a3b8]">60 Questions + Video Solutions • 42 MB</span>
                  </div>
                  <span className="px-2.5 py-1 rounded-lg bg-[#10b981]/20 text-[#10b981] text-xs font-bold border border-[#10b981]/40 flex items-center gap-1">
                    <CheckCircle className="w-3.5 h-3.5" />
                    <span>Saved</span>
                  </span>
                </div>

              </div>

              {/* Data Saving Callout Banner */}
              <div className="p-4 rounded-2xl bg-gradient-to-r from-[#10b981]/15 to-[#059669]/15 border border-[#10b981]/30 flex items-center gap-3">
                <Sparkles className="w-5 h-5 text-[#10b981] shrink-0" />
                <p className="text-xs text-[#d4d4d8]">
                  <strong className="text-white font-semibold">Average Savings:</strong> Students save over <span className="text-[#10b981] font-bold">15 GB</span> of mobile data every month with Mirkuz offline mode.
                </p>
              </div>

            </div>
          </div>

        </div>

      </div>
    </section>
  );
}
