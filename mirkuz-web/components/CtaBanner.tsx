"use client";

import React from "react";
import { ArrowRight, Sparkles, Download, Send } from "lucide-react";

export function CtaBanner() {
  return (
    <section className="bg-[#F5F7FA] py-20 md:py-24 text-center border-b border-[#E8ECF2]">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 space-y-8">
        
        {/* Figma Nexcent Signature Headline */}
        <h2 className="text-3xl sm:text-5xl lg:text-6xl font-bold text-[#263238] tracking-tight leading-[1.15]">
          Pellentesque suscipit <br />
          fringilla libero eu.
        </h2>

        <p className="text-base sm:text-lg text-[#717171] max-w-xl mx-auto">
          Ready to join over 50,000 Ethiopian Grade 12 students and secure your university placement this year?
        </p>

        {/* Green Primary CTA Button */}
        <div className="flex flex-wrap items-center justify-center gap-4 pt-2">
          <a
            href="https://play.google.com/store/apps"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center justify-center gap-2 px-8 py-3.5 rounded-md bg-[#4CAF4F] hover:bg-[#388E3C] text-white font-medium text-base shadow-sm hover:shadow-md hover:-translate-y-0.5 active:translate-y-0 transition-all"
          >
            <span>Get a Demo</span>
            <ArrowRight className="w-4 h-4" />
          </a>

          <a
            href="https://t.me/mirkuz_exam"
            target="_blank"
            rel="noopener noreferrer"
            className="inline-flex items-center justify-center gap-2 px-8 py-3.5 rounded-md bg-white hover:bg-neutral-50 text-[#4D4D4D] border border-[#D5E0D5] font-medium text-base shadow-2xs hover:border-[#4CAF4F] transition-all"
          >
            <Send className="w-4 h-4 text-[#4CAF4F]" />
            <span>Join Telegram</span>
          </a>
        </div>

      </div>
    </section>
  );
}
