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
    <section id="faq" className="bg-white py-20 md:py-28 border-b border-[#E8ECF2]">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Section Header */}
        <div className="text-center space-y-3 mb-14">
          <div className="inline-flex items-center gap-1.5 px-3 py-1 rounded-full bg-[#E8F5E9] text-[#4CAF4F] text-xs font-semibold">
            <HelpCircle className="w-3.5 h-3.5" />
            <span>Got Questions?</span>
          </div>
          <h2 className="text-3xl sm:text-4xl font-bold text-[#4D4D4D] tracking-tight">
            Frequently Asked Questions
          </h2>
          <p className="text-sm sm:text-base text-[#717171]">
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
                className="rounded-lg border border-[#E8ECF2] bg-white overflow-hidden transition-all shadow-[0_2px_4px_rgba(171,190,209,0.15)]"
              >
                <button
                  onClick={() => setOpenIdx(isOpen ? null : idx)}
                  className="w-full p-5 sm:p-6 text-left flex items-center justify-between gap-4 font-semibold text-base sm:text-lg text-[#4D4D4D] hover:text-[#4CAF4F] transition-colors"
                >
                  <span>{faq.q}</span>
                  <ChevronDown
                    className={`w-5 h-5 text-[#89939E] shrink-0 transition-transform duration-200 ${
                      isOpen ? "rotate-180 text-[#4CAF4F]" : ""
                    }`}
                  />
                </button>

                {isOpen && (
                  <div className="px-5 pb-6 sm:px-6 text-sm text-[#717171] leading-relaxed border-t border-[#F5F7FA] pt-4 animate-in fade-in duration-200">
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
