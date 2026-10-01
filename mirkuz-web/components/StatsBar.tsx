"use client";

import React from "react";
import { 
  BookOpen, 
  CheckCircle2, 
  WifiOff, 
  CreditCard, 
  Users, 
  Award,
  Sparkles
} from "lucide-react";

export function StatsBar() {
  const stats = [
    {
      icon: BookOpen,
      value: "50+",
      label: "Complete Subject Courses",
      sub: "Grade 12 Natural & Social streams",
      color: "text-[#10b981]",
      bg: "bg-[#10b981]/10",
      border: "border-[#10b981]/20",
    },
    {
      icon: Award,
      value: "10,000+",
      label: "National Exam Questions",
      sub: "2008 – 2016 E.C. with video solutions",
      color: "text-[#f59e0b]",
      bg: "bg-[#f59e0b]/10",
      border: "border-[#f59e0b]/20",
    },
    {
      icon: WifiOff,
      value: "100%",
      label: "Offline Study Capability",
      sub: "Save videos & PDF notes to phone",
      color: "text-[#38bdf8]",
      bg: "bg-[#38bdf8]/10",
      border: "border-[#38bdf8]/20",
    },
    {
      icon: CreditCard,
      value: "Instant",
      label: "Telebirr & CBE Birr Pay",
      sub: "Automatic activation in seconds",
      color: "text-[#8b5cf6]",
      bg: "bg-[#8b5cf6]/10",
      border: "border-[#8b5cf6]/20",
    },
  ];

  return (
    <section className="relative py-12 bg-[#09090d] border-y border-[#1a1a20]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-6">
          {stats.map((stat, idx) => {
            const Icon = stat.icon;
            return (
              <div
                key={idx}
                className="relative p-6 rounded-2xl bg-[#111115] border border-[#1a1a20] hover:border-[#2c2c35] transition-all hover:translate-y-[-2px] shadow-lg shadow-black/30 group"
              >
                <div className="flex items-center gap-4">
                  <div className={`w-12 h-12 rounded-xl ${stat.bg} ${stat.border} border flex items-center justify-center shrink-0`}>
                    <Icon className={`w-6 h-6 ${stat.color}`} />
                  </div>
                  <div>
                    <div className="text-2xl sm:text-3xl font-extrabold text-white tracking-tight">
                      {stat.value}
                    </div>
                    <div className="text-sm font-semibold text-[#f4f4f5]">
                      {stat.label}
                    </div>
                  </div>
                </div>
                <div className="mt-3 text-xs text-[#94a3b8] font-normal border-t border-[#1a1a20] pt-2">
                  {stat.sub}
                </div>
              </div>
            );
          })}
        </div>
      </div>
    </section>
  );
}
