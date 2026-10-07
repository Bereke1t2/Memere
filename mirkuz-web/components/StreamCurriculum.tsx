"use client";

import React, { useState } from "react";
import { 
  Calculator, 
  Atom, 
  FlaskConical, 
  Dna, 
  BookText, 
  Landmark, 
  Globe2, 
  TrendingUp, 
  CheckCircle2, 
  ArrowRight,
  Sparkles
} from "lucide-react";

export function StreamCurriculum() {
  const [activeStream, setActiveStream] = useState<"natural" | "social">("natural");

  const naturalCourses = [
    {
      title: "Mathematics (Natural)",
      icon: Calculator,
      units: "6 Core Units",
      videos: "48 Video Lessons",
      exams: "2008 - 2016 E.C. Matric Solutions",
      topics: ["Differential Calculus", "Integral Calculus", "Vectors & Matrices", "Sequences & Series", "Probability & Statistics"],
    },
    {
      title: "Physics",
      icon: Atom,
      units: "7 Core Units",
      videos: "42 Video Lessons",
      exams: "2008 - 2016 E.C. Matric Solutions",
      topics: ["Two-Dimensional Motion", "Electromagnetism", "Rotational Dynamics", "Wave Optics & Sound", "Atomic Physics"],
    },
    {
      title: "Chemistry",
      icon: FlaskConical,
      units: "6 Core Units",
      videos: "36 Video Lessons",
      exams: "2008 - 2016 E.C. Matric Solutions",
      topics: ["Chemical Kinetics", "Equilibrium Systems", "Electrochemistry", "Polymers & Organic", "Industrial Chemistry"],
    },
    {
      title: "Biology",
      icon: Dna,
      units: "5 Core Units",
      videos: "32 Video Lessons",
      exams: "2008 - 2016 E.C. Matric Solutions",
      topics: ["Genetics & Evolution", "Human Physiology", "Ecology & Ecosystems", "Microbiology & Biotech", "Cell Biology"],
    },
  ];

  const socialCourses = [
    {
      title: "History of Ethiopia & The World",
      icon: Landmark,
      units: "7 Core Units",
      videos: "40 Video Lessons",
      exams: "2008 - 2016 E.C. Matric Solutions",
      topics: ["Modern Ethiopian State", "Imperial Era & Battles", "World Wars & Cold War", "African Independence", "Historiography"],
    },
    {
      title: "Geography",
      icon: Globe2,
      units: "6 Core Units",
      videos: "34 Video Lessons",
      exams: "2008 - 2016 E.C. Matric Solutions",
      topics: ["GIS & Map Interpretation", "Physical Geography of Ethiopia", "Climatology", "Population Dynamics", "Economic Geography"],
    },
    {
      title: "Economics",
      icon: TrendingUp,
      units: "6 Core Units",
      videos: "38 Video Lessons",
      exams: "2008 - 2016 E.C. Matric Solutions",
      topics: ["Microeconomics & Elasticity", "Macroeconomic Policies", "Fiscal & Monetary Policy", "International Trade", "Development Econ"],
    },
    {
      title: "General Mathematics (Social)",
      icon: Calculator,
      units: "5 Core Units",
      videos: "30 Video Lessons",
      exams: "2008 - 2016 E.C. Matric Solutions",
      topics: ["Commercial Mathematics", "Functions & Graphs", "Statistics & Probability", "Financial Calculations", "Linear Programming"],
    },
  ];

  const currentList = activeStream === "natural" ? naturalCourses : socialCourses;

  return (
    <section id="curriculum" className="bg-[#F8FAFC] py-20 md:py-28 border-b border-[#E2E8F0]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Section Header */}
        <div className="text-center max-w-2xl mx-auto space-y-3 mb-12">
          <div className="inline-flex items-center gap-2 px-3.5 py-1 rounded-full bg-[#E8F8F5] text-[#00B894] text-xs font-semibold">
            <Sparkles className="w-3.5 h-3.5" />
            <span>Ministry Guidelines</span>
          </div>
          <h2 className="text-3xl sm:text-4xl font-extrabold text-[#1E2B58] tracking-tight">
            Comprehensive <span className="title-underline-center text-[#1E2B58]">Stream Curriculum</span>
          </h2>
          <p className="text-sm sm:text-base text-[#64748B] pt-1">
            Aligned 100% with the Ethiopian Ministry of Education matric exam syllabus
          </p>

          {/* Stream Switcher Tab Buttons */}
          <div className="inline-flex p-1.5 rounded-full bg-white border border-[#E2E8F0] shadow-xs mt-6">
            <button
              onClick={() => setActiveStream("natural")}
              className={`px-7 py-2 rounded-full text-xs sm:text-sm font-bold transition-all cursor-pointer ${
                activeStream === "natural"
                  ? "bg-[#1E2B58] text-white shadow-xs"
                  : "text-[#64748B] hover:text-[#1E2B58]"
              }`}
            >
              🌿 Natural Science
            </button>
            <button
              onClick={() => setActiveStream("social")}
              className={`px-7 py-2 rounded-full text-xs sm:text-sm font-bold transition-all cursor-pointer ${
                activeStream === "social"
                  ? "bg-[#1E2B58] text-white shadow-xs"
                  : "text-[#64748B] hover:text-[#1E2B58]"
              }`}
            >
              🏛️ Social Science
            </button>
          </div>
        </div>

        {/* 4 Subject Cards Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
          {currentList.map((course, idx) => {
            const Icon = course.icon;
            return (
              <div
                key={idx}
                className="bg-white rounded-2xl p-6 border border-[#E2E8F0] shadow-dribbble hover:shadow-dribbble-lg hover:-translate-y-1 transition-all flex flex-col justify-between group"
              >
                <div>
                  <div className="w-12 h-12 rounded-xl bg-[#E8F8F5] flex items-center justify-center text-[#00B894] group-hover:scale-110 transition-transform mb-4">
                    <Icon className="w-6 h-6" />
                  </div>

                  <h3 className="text-lg font-bold text-[#1E2B58] mb-2 group-hover:text-[#00B894] transition-colors">
                    {course.title}
                  </h3>

                  <div className="text-xs text-[#64748B] space-y-1 mb-4 pb-3 border-b border-[#F1F5F9]">
                    <div className="font-semibold text-[#00B894]">{course.units} • {course.videos}</div>
                    <div className="text-[#94A3B8]">{course.exams}</div>
                  </div>

                  <div className="space-y-2">
                    <div className="text-[11px] font-bold uppercase tracking-wider text-[#94A3B8]">
                      Key Chapters:
                    </div>
                    {course.topics.map((t, i) => (
                      <div key={i} className="flex items-center gap-2 text-xs text-[#64748B]">
                        <CheckCircle2 className="w-3.5 h-3.5 text-[#00B894] shrink-0" />
                        <span className="truncate">{t}</span>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="mt-6 pt-4 border-t border-[#F1F5F9]">
                  <a
                    href="https://play.google.com/store/apps/details?id=et.mirkuz.mobile"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1.5 text-xs font-bold text-[#00B894] hover:text-[#009874] hover:gap-2.5 transition-all"
                  >
                    <span>Start Studying</span>
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
