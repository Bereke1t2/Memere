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
    { name: "ASTU University", code: "ASTU", icon: Award },
    { name: "Telebirr Ethio Telecom", code: "telebirr", icon: ShieldCheck },
    { name: "Hawassa University", code: "HU", icon: Building2 },
    { name: "Chapa Payment Gateway", code: "chapa", icon: Sparkles },
    { name: "Bahir Dar University", code: "BDU", icon: Globe2 },
    { name: "Jimma University", code: "JU", icon: BookOpen },
  ];

  return (
    <section className="bg-white py-14 border-b border-[#E8ECF2]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
        
        {/* Section Headers */}
        <h2 className="text-2xl sm:text-3xl font-bold text-[#4D4D4D] tracking-tight">
          Our Clients
        </h2>
        <p className="text-sm sm:text-base text-[#717171] mt-2 max-w-xl mx-auto">
          We have been working with some Fortune 500+ clients & top educational institutions
        </p>

        {/* Logos Row */}
        <div className="mt-10 flex flex-wrap items-center justify-center gap-8 sm:gap-12 md:gap-16">
          {clients.map((c, i) => {
            const Icon = c.icon;
            return (
              <div
                key={i}
                className="group flex items-center gap-2 text-[#89939E] hover:text-[#4CAF4F] transition-colors cursor-pointer py-2"
                title={c.name}
              >
                <div className="w-9 h-9 rounded-lg bg-[#F5F7FA] group-hover:bg-[#E8F5E9] flex items-center justify-center transition-colors">
                  <Icon className="w-5 h-5 group-hover:scale-110 transition-transform" />
                </div>
                <span className="font-semibold text-sm tracking-tight text-[#717171] group-hover:text-[#263238] transition-colors">
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
