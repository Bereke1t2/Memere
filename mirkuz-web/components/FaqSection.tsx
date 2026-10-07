"use client";

import React, { useState } from "react";
import { ChevronDown, HelpCircle, Sparkles } from "lucide-react";

export function FaqSection() {
  const [openIdx, setOpenIdx] = useState<number | null>(0);

  const faqs = [
    {
      q: "Does Mirkuz work completely offline without mobile internet?",
      a: "Yes! You can download any video lesson chapter, formula summary sheet, or timed mock exam while connected to Wi-Fi. Once downloaded, you can study, pause, replay, and take exams with mobile data turned off.",
    },
    {
      q: "Are the mock exam questions updated for the 2017 E.C. curriculum?",
      a: "Yes. Our mock exam bank includes authentic past national entrance exams from 2008 to 2016 E.C. plus brand-new predicted mock questions designed by senior Ethiopian examiners specifically aligned with current curriculum guidelines.",
    },
    {
      q: "How do I pay using Telebirr or CBE Birr?",
      a: "Inside the Mirkuz Android app, tap 'Upgrade to All-Access' and choose Telebirr or CBE Birr. You will be redirected to the secure Telebirr checkout or provided a direct USSD/QR payment prompt. Your account activates immediately upon confirmation.",
    },
    {
      q: "Can I use Mirkuz on multiple smartphones or tablets?",
      a: "Your Mirkuz account is tied to your phone number and can be used on your primary Android smartphone or tablet. All your progress, bookmarks, and test scores sync seamlessly.",
    },
    {
      q: "What is included in the Free Starter tier?",
      a: "The Free Starter tier gives you immediate access to introductory video lessons across all 6 subjects, 2 full past entrance exam papers (2015 and 2016 E.C.), and access to our Telegram study community.",
    },
  ];

  return (
    <section id="faq" className="bg-[#F8FAFC] py-20 md:py-28 border-b border-[#E2E8F0]">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Section Header */}
        <div className="text-center space-y-3 mb-14">
          <div className="inline-flex items-center gap-2 px-3.5 py-1 rounded-full bg-[#E8F8F5] text-[#00B894] text-xs font-semibold">
            <HelpCircle className="w-3.5 h-3.5" />
            <span>Got Questions?</span>
          </div>
          <h2 className="text-3xl sm:text-4xl font-extrabold text-[#1E2B58] tracking-tight">
            Frequently Asked <span className="title-underline-center text-[#1E2B58]">Questions</span>
          </h2>
          <p className="text-sm sm:text-base text-[#64748B] pt-1">
            Everything you need to know about the Mirkuz exam prep platform
          </p>
        </div>

        {/* Accordion FAQ Items */}
        <div className="space-y-4">
          {faqs.map((faq, idx) => {
            const isOpen = openIdx === idx;
            return (
              <div
                key={idx}
                className={`rounded-2xl border bg-white overflow-hidden transition-all shadow-xs ${
                  isOpen
                    ? "border-[#00B894]/40 shadow-dribbble"
                    : "border-[#E2E8F0] hover:border-[#CBD5E1]"
                }`}
              >
                <button
                  onClick={() => setOpenIdx(isOpen ? null : idx)}
                  className="w-full p-5 sm:p-6 text-left flex items-center justify-between gap-4 font-bold text-base sm:text-lg text-[#1E2B58] hover:text-[#00B894] transition-colors cursor-pointer"
                >
                  <span>{faq.q}</span>
                  <div className={`w-8 h-8 rounded-full flex items-center justify-center shrink-0 transition-colors ${
                    isOpen ? "bg-[#E8F8F5] text-[#00B894]" : "bg-[#F1F5F9] text-[#94A3B8]"
                  }`}>
                    <ChevronDown
                      className={`w-4 h-4 transition-transform duration-200 ${
                        isOpen ? "rotate-180" : ""
                      }`}
                    />
                  </div>
                </button>

                {isOpen && (
                  <div className="px-5 pb-6 sm:px-6 text-sm text-[#64748B] leading-relaxed border-t border-[#F1F5F9] pt-4 animate-in fade-in duration-200">
                    {faq.a}
                  </div>
                )}
              </div>
            );
          })}
        </div>

      </div>
    </section>
  );
}
