"use client";

import React, { useState } from "react";
import { ChevronDown, HelpCircle, Sparkles } from "lucide-react";

interface FaqItem {
  question: string;
  answer: string;
}

const faqs: FaqItem[] = [
  {
    question: "Does offline mode really work with 0 MB mobile data?",
    answer:
      "Yes! When you are connected to Wi-Fi at home, school, or a cafe, tap the 'Download Unit' button. All HD video lessons, quiz questions, and PDF summaries are securely compressed and saved into your phone's storage. You can then turn off your mobile data completely and study anywhere in Ethiopia without consuming your balance.",
  },
  {
    question: "How do I purchase a course using Telebirr or CBE Birr?",
    answer:
      "When you choose a course or the All-Access Pass inside the Mirkuz app, select 'Telebirr' or 'Chapa (CBE Birr / Bank)'. You will receive an instant 1-click confirmation on your Telebirr app or bank USSD/app. The moment the transaction completes, your courses unlock automatically without any manual code entry.",
  },
  {
    question: "Are Mirkuz courses updated for the new Ethiopian Grade 12 curriculum?",
    answer:
      "Yes. Every course, lesson video, and quiz question has been specifically curated and reviewed according to the latest Ethiopian Ministry of Education (MoE) Grade 12 curriculum standards for both Natural Science and Social Science streams.",
  },
  {
    question: "What years of past National University Entrance Exams are included?",
    answer:
      "Mirkuz includes comprehensive past national entrance exam questions spanning from 2008 E.C. up to the latest 2016 E.C. examinations, complete with step-by-step video solutions and formula derivations by experienced high school teachers.",
  },
  {
    question: "How do I join the Mirkuz Telegram channel?",
    answer:
      "You can tap the 'Join Telegram' button in the navigation bar or bottom banner, or visit t.me/mirkuz_exam. We post daily quiz challenges, entrance exam countdown reminders, and quick formula cheat-sheets for over 15,000 Grade 12 students across Ethiopia.",
  },
  {
    question: "Can I use Mirkuz on multiple devices?",
    answer:
      "Your Mirkuz account can be accessed seamlessly across your Android smartphone and tablet. Your quiz progress, exam attempts, and bookmarks stay synchronized across all your devices.",
  },
];

export function FaqSection() {
  const [openIndex, setOpenIndex] = useState<number | null>(0);

  const toggleAccordion = (idx: number) => {
    setOpenIndex(openIndex === idx ? null : idx);
  };

  return (
    <section id="faq" className="py-24 bg-[#050505] relative">
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Header */}
        <div className="text-center mb-16 space-y-4">
          <div className="inline-flex items-center gap-2 px-3.5 py-1 rounded-full bg-[#10b981]/15 text-[#10b981] border border-[#10b981]/30 text-xs font-bold">
            <HelpCircle className="w-3.5 h-3.5" />
            <span>Got Questions?</span>
          </div>
          <h2 className="text-3xl sm:text-4xl font-extrabold text-white tracking-tight">
            Frequently Asked Questions
          </h2>
          <p className="text-sm sm:text-base text-[#94a3b8]">
            Everything you need to know about preparing for your university entrance exam with Mirkuz.
          </p>
        </div>

        {/* Accordions */}
        <div className="space-y-4">
          {faqs.map((faq, idx) => {
            const isOpen = openIndex === idx;
            return (
              <div
                key={idx}
                className="rounded-2xl bg-[#111115] border border-[#1a1a20] overflow-hidden transition-all"
              >
                <button
                  onClick={() => toggleAccordion(idx)}
                  className="w-full p-6 text-left flex items-center justify-between gap-4 text-white hover:text-[#10b981] transition-colors"
                >
                  <span className="text-base font-bold">{faq.question}</span>
                  <ChevronDown
                    className={`w-5 h-5 text-[#94a3b8] transition-transform duration-200 shrink-0 ${
                      isOpen ? "rotate-180 text-[#10b981]" : ""
                    }`}
                  />
                </button>

                {isOpen && (
                  <div className="px-6 pb-6 pt-1 text-sm text-[#94a3b8] leading-relaxed border-t border-[#1a1a20]/60 animate-in fade-in duration-200">
                    {faq.answer}
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
