"use client";

import React from "react";
import { 
  Users, 
  Building2, 
  Layers, 
  Sparkles,
  ArrowRight,
  WifiOff,
  BookOpen,
  GraduationCap
} from "lucide-react";

export function CommunitySection() {
  const cards = [
    {
      title: "Membership Organisations",
      stream: "Natural Science Stream",
      description:
        "Our membership management software provides full automation of membership renewals and payments with complete Physics, Chemistry, Biology & Math courses.",
      icon: Users,
      badge: "Grade 12 Natural",
    },
    {
      title: "National Associations",
      stream: "Social Science Stream",
      description:
        "Our membership management software provides full automation of membership renewals and payments with in-depth Geography, History, Economics & Aptitude.",
      icon: Building2,
      badge: "Grade 12 Social",
    },
    {
      title: "Clubs And Groups",
      stream: "100% Offline Learning",
      description:
        "Our membership management software provides full automation of membership renewals and payments with encrypted offline downloads and zero data usage.",
      icon: Layers,
      badge: "Zero Data",
    },
  ];

  return (
    <section id="services" className="bg-white py-20 border-b border-[#E8ECF2]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Section Header */}
        <div className="text-center max-w-2xl mx-auto space-y-2 mb-14">
          <h2 className="text-3xl sm:text-4xl font-bold text-[#4D4D4D] tracking-tight">
            Manage your entire community <br className="hidden sm:inline" />
            in a single system
          </h2>
          <p className="text-sm sm:text-base text-[#717171]">
            Who is Nextcent suitable for?
          </p>
        </div>

        {/* 3-Card Grid */}
        <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
          {cards.map((card, idx) => {
            const Icon = card.icon;
            return (
              <div
                key={idx}
                className="group relative bg-white rounded-lg p-8 text-center border border-[#E8ECF2] shadow-[0_2px_4px_rgba(171,190,209,0.2)] hover:shadow-[0_8px_16px_rgba(171,190,209,0.3)] hover:-translate-y-1 transition-all flex flex-col justify-between"
              >
                <div>
                  {/* Icon Circle */}
                  <div className="relative w-16 h-16 mx-auto mb-6 flex items-center justify-center">
                    <div className="w-14 h-14 rounded-2xl bg-[#E8F5E9] group-hover:bg-[#4CAF4F] flex items-center justify-center transition-colors">
                      <Icon className="w-7 h-7 text-[#4CAF4F] group-hover:text-white transition-colors" />
                    </div>
                  </div>

                  {/* Title */}
                  <h3 className="text-xl sm:text-2xl font-bold text-[#4D4D4D] mb-3 leading-snug">
                    {card.title}
                  </h3>

                  {/* Stream Subtitle Badge */}
                  <div className="inline-block px-2.5 py-0.5 mb-3 rounded-full bg-[#F5F7FA] text-[#4CAF4F] text-xs font-semibold">
                    {card.stream}
                  </div>

                  {/* Description */}
                  <p className="text-sm text-[#717171] leading-relaxed">
                    {card.description}
                  </p>
                </div>

                {/* Subtle Learn More Link */}
                <div className="mt-6 pt-4 border-t border-[#F5F7FA]">
                  <a
                    href="#curriculum"
                    className="inline-flex items-center gap-1.5 text-xs font-semibold text-[#4CAF4F] hover:text-[#388E3C] transition-colors"
                  >
                    <span>View Curriculum</span>
                    <ArrowRight className="w-3.5 h-3.5" />
                  </a>
                </div>
              </div>
            );
          })}
        </div>

      </div>
    </section>
  );
}
