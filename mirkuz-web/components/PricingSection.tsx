"use client";

import React from "react";
import { 
  Check, 
  Sparkles, 
  CreditCard, 
  ShieldCheck, 
  Zap, 
  Download,
  Smartphone
} from "lucide-react";

export function PricingSection() {
  const plans = [
    {
      name: "Free Starter",
      price: "0",
      currency: "ETB",
      period: "Forever Free",
      description: "Sample core video lessons and try past national entrance exam questions.",
      badge: null,
      highlight: false,
      features: [
        "Sample introductory video lessons in all 6 subjects",
        "2 Past National Exam Papers (2015 & 2016 E.C.)",
        "Basic downloadable PDF chapter summaries",
        "In-app question practice mode",
        "Community Telegram study access",
      ],
      ctaText: "Get Started Free",
      ctaLink: "https://play.google.com/store/apps",
    },
    {
      name: "All-Access Entrance Pass",
      price: "1,450",
      currency: "ETB",
      period: "Full Academic Year",
      description: "Everything you need to score 550+ on the Ethiopian University Entrance Exam.",
      badge: "MOST POPULAR",
      highlight: true,
      features: [
        "Complete video courses for ALL 6 subjects (Natural or Social)",
        "10,000+ Past National Exam Questions (2008 – 2016 E.C.)",
        "Step-by-step video solutions for every past exam question",
        "Full timed mock exams with real countdown timers",
        "100% Unlimited encrypted offline downloads",
        "Comprehensive PDF chapter study notes & formula sheets",
        "Mirkuz AI Mascot concept tutor assistance",
        "Diagnostic weak-area score improvement analytics",
      ],
      ctaText: "Unlock All-Access on Google Play",
      ctaLink: "https://play.google.com/store/apps",
    },
    {
      name: "Single Subject Pass",
      price: "350",
      currency: "ETB",
      period: "Per Subject / Year",
      description: "Focus on mastering a specific difficult subject like Mathematics or Physics.",
      badge: "FLEXIBLE",
      highlight: false,
      features: [
        "Complete video library for 1 selected subject",
        "All chapter quizzes & unit assessments",
        "10-Year past national exam question bank for that subject",
        "Full downloadable unit PDF notes",
        "Offline video saving for that subject",
        "Subject-specific progress tracking",
      ],
      ctaText: "Select Subject in App",
      ctaLink: "https://play.google.com/store/apps",
    },
  ];

  return (
    <section id="pricing" className="py-24 bg-[#050505] relative overflow-hidden">
      {/* Background radial glow */}
      <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-[800px] h-[500px] bg-[#10b981]/10 rounded-full blur-[160px] pointer-events-none" />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        
        {/* Header */}
        <div className="text-center max-w-3xl mx-auto mb-16 space-y-4">
          <div className="inline-flex items-center gap-2 px-3.5 py-1 rounded-full bg-[#10b981]/15 text-[#10b981] border border-[#10b981]/30 text-xs font-bold">
            <CreditCard className="w-3.5 h-3.5" />
            <span>Transparent Ethiopian Birr (ETB) Pricing</span>
          </div>
          <h2 className="text-3xl sm:text-4xl font-extrabold text-white tracking-tight">
            Invest in Your University Admission at an Affordable Price
          </h2>
          <p className="text-sm sm:text-base text-[#94a3b8]">
            Pay seamlessly with Telebirr or CBE Birr with instant automated activation. No recurring credit card hassle.
          </p>
        </div>

        {/* Pricing Cards */}
        <div className="grid grid-cols-1 lg:grid-cols-3 gap-8 items-stretch">
          {plans.map((plan, idx) => (
            <div
              key={idx}
              className={`rounded-3xl p-8 flex flex-col justify-between transition-all relative ${
                plan.highlight
                  ? "bg-[#111115] border-2 border-[#10b981] shadow-2xl shadow-[#10b981]/20 scale-105 z-20"
                  : "bg-[#111115]/90 border border-[#1a1a20] hover:border-[#2c2c35]"
              }`}
            >
              {plan.badge && (
                <div className="absolute -top-3.5 left-1/2 -translate-x-1/2 px-4 py-1 rounded-full bg-gradient-to-r from-[#10b981] to-[#059669] text-white text-[11px] font-bold tracking-wider uppercase shadow-lg shadow-[#10b981]/30">
                  {plan.badge}
                </div>
              )}

              <div>
                {/* Plan Title & Price */}
                <div className="mb-6 space-y-2">
                  <h3 className="text-xl font-bold text-white">{plan.name}</h3>
                  <p className="text-xs text-[#94a3b8]">{plan.description}</p>
                  
                  <div className="pt-4 flex items-baseline gap-1">
                    <span className="text-3xl sm:text-4xl font-extrabold text-white">
                      {plan.price}
                    </span>
                    <span className="text-base font-bold text-[#10b981]">{plan.currency}</span>
                    <span className="text-xs text-[#94a3b8] ml-2 font-medium">/ {plan.period}</span>
                  </div>
                </div>

                {/* Features List */}
                <div className="space-y-3 pt-6 border-t border-[#1a1a20]">
                  <span className="text-xs font-semibold text-[#d4d4d8] uppercase tracking-wider block">
                    What&apos;s Included:
                  </span>
                  <ul className="space-y-2.5">
                    {plan.features.map((feat, fIdx) => (
                      <li key={fIdx} className="flex items-start gap-2.5 text-xs text-[#d4d4d8] leading-relaxed">
                        <Check className="w-4 h-4 text-[#10b981] shrink-0 mt-0.5" />
                        <span>{feat}</span>
                      </li>
                    ))}
                  </ul>
                </div>
              </div>

              {/* Action Button */}
              <div className="mt-8 pt-6 border-t border-[#1a1a20]">
                <a
                  href={plan.ctaLink}
                  target="_blank"
                  rel="noopener noreferrer"
                  className={`w-full py-3.5 px-4 rounded-xl text-xs font-bold flex items-center justify-center gap-2 transition-all ${
                    plan.highlight
                      ? "bg-gradient-to-r from-[#10b981] to-[#059669] text-white hover:opacity-95 shadow-lg shadow-[#10b981]/25"
                      : "bg-[#18181e] text-white hover:bg-[#25252e] border border-[#2c2c35]"
                  }`}
                >
                  <Download className="w-4 h-4" />
                  <span>{plan.ctaText}</span>
                </a>
              </div>
            </div>
          ))}
        </div>

        {/* Local Payment Methods Banner */}
        <div className="mt-16 p-6 rounded-2xl bg-[#111115] border border-[#1a1a20] flex flex-wrap items-center justify-between gap-6">
          <div className="flex items-center gap-3">
            <div className="w-10 h-10 rounded-xl bg-[#f59e0b]/15 flex items-center justify-center text-[#f59e0b]">
              <ShieldCheck className="w-5 h-5" />
            </div>
            <div>
              <h4 className="text-sm font-bold text-white">Instant Local Payment Activation</h4>
              <p className="text-xs text-[#94a3b8]">Courses unlock automatically the moment payment is completed.</p>
            </div>
          </div>

          <div className="flex items-center gap-3 flex-wrap">
            <span className="px-3 py-1 rounded-lg bg-[#18181e] border border-[#2c2c35] text-xs font-bold text-[#f59e0b]">
              📱 Telebirr (Ethio Telecom)
            </span>
            <span className="px-3 py-1 rounded-lg bg-[#18181e] border border-[#2c2c35] text-xs font-bold text-[#38bdf8]">
              🏦 CBE Birr
            </span>
            <span className="px-3 py-1 rounded-lg bg-[#18181e] border border-[#2c2c35] text-xs font-bold text-[#10b981]">
              💳 Chapa (All Banks & Cards)
            </span>
          </div>
        </div>

      </div>
    </section>
  );
}
