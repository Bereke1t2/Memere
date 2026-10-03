"use client";

import React from "react";
import { 
  Users, 
  Building2, 
  MousePointerClick, 
  CreditCard,
  BookOpen,
  GraduationCap,
  Award,
  CheckCircle2
} from "lucide-react";

export function StatsCounter() {
  const stats = [
    {
      icon: Users,
      value: "2,245,341",
      label: "Members",
      sub: "Active Grade 12 Students",
    },
    {
      icon: Building2,
      value: "46,328",
      label: "Clubs",
      sub: "High Schools & Study Groups",
    },
    {
      icon: MousePointerClick,
      value: "828,867",
      label: "Event Bookings",
      sub: "Mock Exam Tests Taken",
    },
    {
      icon: CreditCard,
      value: "1,926,436",
      label: "Payments",
      sub: "Telebirr & Local Gateways",
    },
  ];

  return (
    <section className="bg-[#F5F7FA] py-16 sm:py-20 border-b border-[#E8ECF2]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-10 lg:gap-12 items-center">
          
          {/* Left Column: Heading */}
          <div className="lg:col-span-5 space-y-3">
            <h2 className="text-3xl sm:text-4xl font-bold text-[#4D4D4D] tracking-tight leading-tight">
              Helping a local <br />
              <span className="text-[#4CAF4F]">business reinvent itself</span>
            </h2>
            <p className="text-sm sm:text-base text-[#717171]">
              We reached here with our hard work and dedication
            </p>
          </div>

          {/* Right Column: 2x2 Metric Grid */}
          <div className="lg:col-span-7 grid grid-cols-1 sm:grid-cols-2 gap-8">
            {stats.map((stat, idx) => {
              const Icon = stat.icon;
              return (
                <div key={idx} className="flex items-center gap-4 group">
                  <div className="w-12 h-12 rounded-xl bg-white group-hover:bg-[#E8F5E9] border border-[#E8ECF2] flex items-center justify-center text-[#4CAF4F] shadow-xs transition-colors shrink-0">
                    <Icon className="w-6 h-6" />
                  </div>
                  <div>
                    <div className="text-2xl sm:text-3xl font-bold text-[#4D4D4D] tracking-tight">
                      {stat.value}
                    </div>
                    <div className="text-sm font-medium text-[#717171]">
                      {stat.label}
                    </div>
                    <div className="text-xs text-[#89939E]">
                      {stat.sub}
                    </div>
                  </div>
                </div>
              );
            })}
          </div>

        </div>
      </div>
    </section>
  );
}
