"use client";

import React from "react";
import Image from "next/image";
import { 
  ArrowRight, 
  Building2, 
  GraduationCap, 
  ShieldCheck, 
  Sparkles,
  Award,
  Globe2,
  BookOpen
} from "lucide-react";

export function TestimonialSpotlight() {
  const partnerIcons = [
    { code: "AAU", icon: GraduationCap },
    { code: "ASTU", icon: Award },
    { code: "telebirr", icon: ShieldCheck },
    { code: "HU", icon: Building2 },
    { code: "chapa", icon: Sparkles },
    { code: "BDU", icon: Globe2 },
  ];

  return (
    <section id="testimonial" className="bg-[#F5F7FA] py-20 border-b border-[#E8ECF2]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-10 lg:gap-12 items-center">
          
          {/* Left Column: Portrait Avatar / Feature Card */}
          <div className="lg:col-span-4 flex justify-center">
            <div className="relative w-full max-w-[320px] aspect-square rounded-2xl overflow-hidden shadow-md border border-[#E8ECF2] bg-white flex flex-col items-center justify-center p-6 text-center">
              <div className="w-28 h-28 rounded-full bg-[#E8F5E9] border-4 border-white shadow-sm flex items-center justify-center text-[#4CAF4F] mb-4">
                <GraduationCap className="w-14 h-14" />
              </div>
              <h4 className="text-lg font-bold text-[#4D4D4D]">Kalkidan Bekele</h4>
              <p className="text-xs font-semibold text-[#4CAF4F] mt-0.5">Scored 642 / 700 (2016 E.C.)</p>
              <p className="text-xs text-[#717171] mt-1">Addis Ababa University — Medicine</p>
            </div>
          </div>

          {/* Right Column: Quote & Author & Partners */}
          <div className="lg:col-span-8 flex flex-col space-y-5">
            
            <p className="text-base sm:text-lg text-[#717171] leading-relaxed italic">
              “Maecenas dignissim justo eget nulla rutrum molestie. Maecenas lobortis sem dui, vel rutrum risus tincidunt ullamcorper. Proin eu enim metus. Vivamus sed libero ornare, tristique quam in, gravida enim. Nullam ut molestie arcu, at hendrerit elit. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Morbi eu quam non.”
            </p>

            <div>
              <div className="text-lg font-bold text-[#4CAF4F]">
                Tim Smith
              </div>
              <div className="text-sm text-[#89939E]">
                British Dragon Boat Racing Club / AAU Medical Faculty
              </div>
            </div>

            {/* Bottom Row: Client Logos + "Meet all customers ->" Link */}
            <div className="pt-4 flex flex-wrap items-center justify-between gap-6 border-t border-[#E8ECF2]/60">
              
              {/* Partner Icons */}
              <div className="flex flex-wrap items-center gap-6 text-[#89939E]">
                {partnerIcons.map((p, i) => {
                  const Icon = p.icon;
                  return (
                    <div key={i} className="flex items-center gap-1.5 hover:text-[#4CAF4F] transition-colors cursor-pointer">
                      <Icon className="w-4 h-4" />
                      <span className="text-xs font-medium">{p.code}</span>
                    </div>
                  );
                })}
              </div>

              {/* Meet all customers link */}
              <a
                href="#testimonials"
                className="inline-flex items-center gap-2 text-sm sm:text-base font-bold text-[#4CAF4F] hover:text-[#388E3C] hover:gap-3 transition-all"
              >
                <span>Meet all customers</span>
                <ArrowRight className="w-4 h-4" />
              </a>

            </div>

          </div>

        </div>
      </div>
    </section>
  );
}
