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
    <section id="curriculum" className="bg-white py-20 border-b border-[#E8ECF2]">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Section Header */}
        <div className="text-center max-w-2xl mx-auto space-y-3 mb-10">
          <h2 className="text-3xl sm:text-4xl font-bold text-[#4D4D4D] tracking-tight">
            Comprehensive Stream Curriculum
          </h2>
          <p className="text-sm sm:text-base text-[#717171]">
            Aligned 100% with the Ethiopian Ministry of Education matrix exam guidelines
          </p>

          {/* Stream Switcher Tab Buttons */}
          <div className="inline-flex p-1 rounded-lg bg-[#F5F7FA] border border-[#E8ECF2] shadow-2xs mt-4">
            <button
              onClick={() => setActiveStream("natural")}
              className={`px-6 py-2.5 rounded-md text-xs sm:text-sm font-bold transition-all ${
                activeStream === "natural"
                  ? "bg-[#4CAF4F] text-white shadow-xs"
                  : "text-[#717171] hover:text-[#4D4D4D]"
              }`}
            >
              🌿 Natural Science
            </button>
            <button
              onClick={() => setActiveStream("social")}
              className={`px-6 py-2.5 rounded-md text-xs sm:text-sm font-bold transition-all ${
                activeStream === "social"
                  ? "bg-[#4CAF4F] text-white shadow-xs"
                  : "text-[#717171] hover:text-[#4D4D4D]"
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
                className="bg-white rounded-lg p-6 border border-[#E8ECF2] shadow-[0_2px_4px_rgba(171,190,209,0.2)] hover:shadow-[0_8px_16px_rgba(171,190,209,0.3)] hover:-translate-y-1 transition-all flex flex-col justify-between"
              >
                <div>
                  <div className="w-12 h-12 rounded-xl bg-[#E8F5E9] flex items-center justify-center text-[#4CAF4F] mb-4">
                    <Icon className="w-6 h-6" />
                  </div>

                  <h3 className="text-lg font-bold text-[#4D4D4D] mb-2">
                    {course.title}
                  </h3>

                  <div className="text-xs text-[#717171] space-y-1 mb-4 pb-3 border-b border-[#F5F7FA]">
                    <div className="font-semibold text-[#4CAF4F]">{course.units} • {course.videos}</div>
                    <div className="text-[#89939E]">{course.exams}</div>
                  </div>

                  <div className="space-y-1.5">
                    <div className="text-[11px] font-bold uppercase tracking-wider text-[#89939E]">
                      Key Chapters:
                    </div>
                    {course.topics.map((t, i) => (
                      <div key={i} className="flex items-center gap-1.5 text-xs text-[#717171]">
                        <CheckCircle2 className="w-3.5 h-3.5 text-[#4CAF4F] shrink-0" />
                        <span className="truncate">{t}</span>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="mt-6 pt-4 border-t border-[#F5F7FA]">
                  <a
                    href="https://play.google.com/store/apps"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="inline-flex items-center gap-1.5 text-xs font-bold text-[#4CAF4F] hover:text-[#388E3C] transition-colors"
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
