"use client";

import React from "react";
import { 
  Check, 
  Sparkles, 
  CreditCard, 
  ShieldCheck, 
  Zap, 
  Download,
  Smartphone,
  ArrowRight
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
        "Sample introductory video lessons in all subjects",
        "2 Past National Exam Papers (2015 & 2016 E.C.)",
        "Basic downloadable PDF chapter summaries",
        "In-app question practice mode",
        "Community Telegram study access",
      ],
      ctaText: "Get Started Free",
      ctaLink: "https://play.google.com/store/apps/details?id=et.mirkuz.mobile",
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
        "All 50+ Grade 12 Video Masterclasses (HD)",
        "10,000+ Past Matric Questions (2008 – 2016 E.C.)",
        "Step-by-step Examiner Video Explanations",
        "100% Data-Free Offline Encrypted Storage",
        "Full Timed National Mock Exam Simulator",
        "Direct Telebirr & CBE Birr Instant Activation",
      ],
      ctaText: "Unlock All-Access Pass",
      ctaLink: "https://play.google.com/store/apps/details?id=et.mirkuz.mobile",
    },
    {
      name: "School & Group Pass",
      price: "950",
      currency: "ETB / Student",
      period: "Min 10 Students",
      description: "Discounted bulk access for high schools, tutoring centers, and study clubs.",
      badge: "GROUP SAVINGS",
      highlight: false,
      features: [
        "Everything in the All-Access Pass included",
        "Teacher & School Administrator Dashboard",
        "Student class analytics & weak-area diagnostics",
        "Bulk Telebirr / CBE invoice payments",
        "Dedicated Telegram priority support manager",
      ],
      ctaText: "Contact for School Pass",
      ctaLink: "https://t.me/mirkuz_exam",
    },
  ];

  return (
    <section id="pricing" className="bg-white py-20 md:py-28 border-b border-[#E2E8F0]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Section Header */}
        <div className="text-center max-w-2xl mx-auto space-y-3 mb-16">
          <div className="inline-flex items-center gap-2 px-3.5 py-1 rounded-full bg-[#E8F8F5] text-[#00B894] text-xs font-semibold">
            <ShieldCheck className="w-3.5 h-3.5" />
            <span>Fair & Transparent Local Pricing</span>
          </div>
          <h2 className="text-3xl sm:text-4xl font-extrabold text-[#1E2B58] tracking-tight">
            Simple, Transparent <span className="title-underline-center text-[#1E2B58]">ETB Pricing</span>
          </h2>
          <p className="text-sm sm:text-base text-[#64748B] pt-1">
            One-time academic year fee. Pay seamlessly with Telebirr, CBE Birr, or Chapa. Zero recurring surprise charges.
          </p>
        </div>

        {/* Pricing Cards Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8 items-stretch">
          {plans.map((plan, idx) => (
            <div
              key={idx}
              className={`rounded-3xl p-8 flex flex-col justify-between transition-all relative ${
                plan.highlight
                  ? "bg-white border-2 border-[#00B894] shadow-dribbble-lg md:-translate-y-2 ring-4 ring-[#00B894]/10"
                  : "bg-white border border-[#E2E8F0] shadow-dribbble hover:shadow-dribbble-lg"
              }`}
            >
              {/* Badge for Popular Plan */}
              {plan.badge && (
                <div className="absolute -top-3.5 left-1/2 -translate-x-1/2 px-4 py-1 rounded-full bg-[#00B894] text-white text-[10px] font-extrabold tracking-wider uppercase shadow-xs">
                  {plan.badge}
                </div>
              )}

              <div>
                <h3 className="text-xl font-bold text-[#1E2B58] mb-1">
                  {plan.name}
                </h3>
                <p className="text-xs text-[#64748B] mb-6 min-h-[36px] leading-relaxed">
                  {plan.description}
                </p>

                {/* Price Display */}
                <div className="flex items-baseline gap-1.5 pb-6 mb-6 border-b border-[#F1F5F9]">
                  <span className="text-4xl font-extrabold text-[#1E2B58] tracking-tight">
                    {plan.price}
                  </span>
                  <span className="text-sm font-bold text-[#00B894]">
                    {plan.currency}
                  </span>
                  <span className="text-xs text-[#94A3B8] ml-1">
                    / {plan.period}
                  </span>
                </div>

                {/* Feature List */}
                <div className="space-y-3 mb-8">
                  {plan.features.map((feat, i) => (
                    <div key={i} className="flex items-start gap-2.5 text-xs text-[#1E2B58]">
                      <div className="w-4 h-4 rounded-full bg-[#E8F8F5] flex items-center justify-center text-[#00B894] shrink-0 mt-0.5">
                        <Check className="w-3 h-3 stroke-[3]" />
                      </div>
                      <span className="leading-snug text-[#475569]">{feat}</span>
                    </div>
                  ))}
                </div>
              </div>

              {/* CTA Button */}
              <a
                href={plan.ctaLink}
                target="_blank"
                rel="noopener noreferrer"
                className={`w-full py-3.5 px-6 rounded-full text-xs font-bold flex items-center justify-center gap-2 transition-all cursor-pointer ${
                  plan.highlight
                    ? "bg-[#1E2B58] hover:bg-[#152042] text-white shadow-md hover:shadow-lg"
                    : "bg-[#F1F5F9] hover:bg-[#E8F8F5] text-[#1E2B58] hover:text-[#00B894] border border-[#E2E8F0]"
                }`}
              >
                <span>{plan.ctaText}</span>
                <ArrowRight className="w-3.5 h-3.5" />
              </a>

            </div>
          ))}
        </div>

      </div>
    </section>
  );
}
