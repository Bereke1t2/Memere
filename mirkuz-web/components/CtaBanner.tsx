"use client";

import React from "react";
import { Download, Send, ArrowRight, Sparkles, QrCode } from "lucide-react";

export function CtaBanner() {
  return (
    <section className="py-20 bg-[#09090d] relative overflow-hidden">
      {/* Glow */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[700px] h-[350px] bg-gradient-to-r from-[#10b981]/20 via-[#f59e0b]/15 to-transparent rounded-full blur-[140px] pointer-events-none" />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        <div className="rounded-3xl bg-gradient-to-br from-[#18181e] via-[#111115] to-[#0a0a0d] border border-[#2c2c35] p-8 sm:p-12 lg:p-16 shadow-2xl relative overflow-hidden">
          
          <div className="grid grid-cols-1 lg:grid-cols-12 gap-8 items-center">
            
            {/* Left Column: Heading & CTAs */}
            <div className="lg:col-span-8 space-y-6 text-left">
              <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-[#10b981]/15 text-[#10b981] border border-[#10b981]/30 text-xs font-bold">
                <Sparkles className="w-3.5 h-3.5" />
                <span>Join Over 15,000 Ethiopian Grade 12 Students</span>
              </div>

              <h2 className="text-3xl sm:text-4xl lg:text-5xl font-extrabold text-white tracking-tight leading-tight">
                Ready to Ace Your University Entrance Exam?
              </h2>

              <p className="text-base text-[#94a3b8] max-w-xl leading-relaxed">
                Download Mirkuz on Google Play today, or join our Telegram channel for daily entrance exam practice and study tips.
              </p>

              <div className="flex flex-col sm:flex-row items-stretch sm:items-center gap-4 pt-2">
                <a
                  href="https://play.google.com/store/apps"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center justify-center gap-3 px-6 py-4 rounded-xl bg-gradient-to-r from-[#10b981] to-[#059669] text-white font-bold shadow-xl shadow-[#10b981]/25 hover:shadow-2xl hover:shadow-[#10b981]/40 hover:scale-[1.02] transition-all"
                >
                  <Download className="w-5 h-5" />
                  <span>Get Mirkuz on Google Play</span>
                  <ArrowRight className="w-4 h-4 ml-1" />
                </a>

                <a
                  href="https://t.me/mirkuz_exam"
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center justify-center gap-3 px-6 py-4 rounded-xl bg-[#18181e] hover:bg-[#25252e] text-white font-bold border border-[#2c2c35] hover:border-[#38bdf8]/50 transition-all"
                >
                  <Send className="w-5 h-5 text-[#38bdf8]" />
                  <span>Join Telegram Community</span>
                </a>
              </div>
            </div>

            {/* Right Column: QR Code & Mobile Fast Access */}
            <div className="lg:col-span-4 flex justify-center lg:justify-end">
              <div className="p-6 rounded-2xl bg-[#09090b] border border-[#2c2c35] text-center space-y-3 max-w-[220px]">
                <div className="w-32 h-32 mx-auto rounded-xl bg-white p-2.5 flex items-center justify-center shadow-lg">
                  {/* Simulated SVG QR Code */}
                  <svg className="w-full h-full text-black" viewBox="0 0 24 24" fill="currentColor">
                    <path d="M2 2h8v8H2V2zm2 2v4h4V4H4zm10-2h8v8h-8V2zm2 2v4h4V4h-4zM2 14h8v8H2v-8zm2 2v4h4v-4H4zm14 0h4v4h-4v-4zm-4 4h2v2h-2v-2zm0-4h2v2h-2v-2zm4-4h2v2h-2v-2zm-4 0h2v2h-2v-2zm2 2h2v2h-2v-2zM5 5h2v2H5V5zm12 0h2v2h-2V5zM5 17h2v2H5v-2z" />
                  </svg>
                </div>
                <div>
                  <span className="text-xs font-bold text-white block">Scan to Install</span>
                  <span className="text-[10px] text-[#94a3b8]">Android APK / Google Play</span>
                </div>
              </div>
            </div>

          </div>

        </div>
      </div>
    </section>
  );
}
