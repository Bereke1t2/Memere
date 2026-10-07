"use client";

import React from "react";
import { 
  Building2, 
  GraduationCap, 
  ShieldCheck, 
  Sparkles,
  Award,
  Globe2,
  BookOpen
} from "lucide-react";

export function ClientsBar() {
  const clients = [
    { name: "Addis Ababa University", code: "AAU", icon: GraduationCap },
    { name: "Addis Ababa Science & Tech University", code: "AASTU", icon: Award },
    { name: "Adama Science & Tech University", code: "ASTU", icon: Building2 },
    { name: "Jimma University", code: "JU", icon: BookOpen },
    { name: "Hawassa University", code: "HU", icon: Globe2 },
    { name: "Bahir Dar University", code: "BDU", icon: GraduationCap },
    { name: "Telebirr Mobile Money", code: "Telebirr", icon: ShieldCheck },
    { name: "Chapa Payment Gateway", code: "Chapa", icon: Sparkles },
  ];

  return (
    <section className="bg-white py-14 border-y border-[#F1F5F9]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        
        {/* Section Headers with Dribbble Accent Underline */}
        <h2 className="text-2xl sm:text-3xl font-extrabold text-[#1E2B58] tracking-tight">
          Admissions & Partner Institutions
        </h2>
        <div className="title-underline-center" />

        <p className="text-xs sm:text-sm text-[#64748B] mt-3 max-w-xl mx-auto">
          Preparing students for competitive placement into Ethiopia's top public universities and technology institutes
        </p>

        {/* Logos Row with Dribbble Elevated Chips */}
        <div className="mt-9 flex flex-wrap items-center justify-center gap-4 sm:gap-6 md:gap-8">
          {clients.map((c, i) => {
            const Icon = c.icon;
            return (
              <div
                key={i}
                className="group flex items-center gap-2.5 px-4 py-2.5 rounded-2xl bg-white border border-[#E2E8F0] shadow-sm hover:border-[#00B894] hover:shadow-dribbble transition-all cursor-pointer hover:-translate-y-0.5"
                title={c.name}
              >
                <div className="w-8 h-8 rounded-xl bg-[#E8F8F5] group-hover:bg-[#00B894] flex items-center justify-center text-[#00B894] group-hover:text-white transition-colors">
                  <Icon className="w-4 h-4 transition-transform group-hover:scale-110" />
                </div>
                <span className="font-bold text-xs tracking-tight text-[#1E2B58] group-hover:text-[#00B894] transition-colors">
                  {c.code}
                </span>
              </div>
            );
          })}
        </div>

      </div>
    </section>
  );
}
