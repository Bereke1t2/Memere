"use client";

import React, { useState } from "react";
import { 
  Atom, 
  FlaskConical, 
  Dna, 
  Calculator, 
  BookText, 
  Landmark, 
  Globe2, 
  TrendingUp, 
  FileText, 
  PlayCircle, 
  CheckCircle2, 
  Sparkles,
  Layers
} from "lucide-react";

export function StreamCurriculum() {
  const [activeStream, setActiveStream] = useState<"natural" | "social">("natural");

  const naturalCourses = [
    {
      title: "Mathematics (Natural)",
      icon: Calculator,
      color: "text-[#10b981]",
      badgeColor: "bg-[#10b981]/15 text-[#10b981] border-[#10b981]/30",
      units: "6 Core Units",
      videos: "48 Video Lessons (24 hrs)",
      exams: "2008 - 2016 E.C. Solutions",
      topics: ["Sequences & Series", "Vectors & Matrices", "Differential Calculus", "Integral Calculus", "Coordinate Geometry", "Probability"],
    },
    {
      title: "Physics",
      icon: Atom,
      color: "text-[#38bdf8]",
      badgeColor: "bg-[#38bdf8]/15 text-[#38bdf8] border-[#38bdf8]/30",
      units: "7 Core Units",
      videos: "42 Video Lessons (20 hrs)",
      exams: "2008 - 2016 E.C. Solutions",
      topics: ["Thermodynamics", "Electromagnetism", "Optics & Wave Motion", "Atomic & Nuclear Physics", "Fluid Dynamics", "Rotational Dynamics"],
    },
    {
      title: "Chemistry",
      icon: FlaskConical,
      color: "text-[#f59e0b]",
      badgeColor: "bg-[#f59e0b]/15 text-[#f59e0b] border-[#f59e0b]/30",
      units: "6 Core Units",
      videos: "38 Video Lessons (18 hrs)",
      exams: "2008 - 2016 E.C. Solutions",
      topics: ["Reaction Kinetics", "Chemical Equilibrium", "Acid-Base Equilibria", "Electrochemistry", "Polymers & Organic Chemistry", "Industrial Chemistry"],
    },
    {
      title: "Biology",
      icon: Dna,
      color: "text-[#ec4899]",
      badgeColor: "bg-[#ec4899]/15 text-[#ec4899] border-[#ec4899]/30",
      units: "5 Core Units",
      videos: "34 Video Lessons (16 hrs)",
      exams: "2008 - 2016 E.C. Solutions",
      topics: ["Molecular Genetics", "Human Biology & Physiology", "Cellular Respiration & Photosynthesis", "Ecology & Conservation", "Biotechnology"],
    },
    {
      title: "English (Advanced)",
      icon: BookText,
      color: "text-[#8b5cf6]",
      badgeColor: "bg-[#8b5cf6]/15 text-[#8b5cf6] border-[#8b5cf6]/30",
      units: "5 Core Units",
      videos: "30 Video Lessons (14 hrs)",
      exams: "2008 - 2016 E.C. Solutions",
      topics: ["Reading Comprehension Strategies", "Advanced Grammar & Mechanics", "Vocabulary in Context", "Sentence Structure", "Exam Essay Techniques"],
    },
    {
      title: "Civics & General Aptitude",
      icon: Landmark,
      color: "text-[#14b8a6]",
      badgeColor: "bg-[#14b8a6]/15 text-[#14b8a6] border-[#14b8a6]/30",
      units: "5 Core Units",
      videos: "26 Video Lessons (12 hrs)",
      exams: "2008 - 2016 E.C. Solutions",
      topics: ["Constitution & Federalism", "Rule of Law & Human Rights", "Logical Reasoning", "Quantitative Aptitude", "Analytical Problem Solving"],
    },
  ];

  const socialCourses = [
    {
      title: "General Mathematics (Social)",
      icon: Calculator,
      color: "text-[#10b981]",
      badgeColor: "bg-[#10b981]/15 text-[#10b981] border-[#10b981]/30",
      units: "5 Core Units",
      videos: "36 Video Lessons (18 hrs)",
      exams: "2008 - 2016 E.C. Solutions",
      topics: ["Commercial Mathematics", "Business Statistics", "Probability", "Linear Programming", "Functions & Graphs"],
    },
    {
      title: "History of Ethiopia & The World",
      icon: Landmark,
      color: "text-[#f59e0b]",
      badgeColor: "bg-[#f59e0b]/15 text-[#f59e0b] border-[#f59e0b]/30",
      units: "6 Core Units",
      videos: "40 Video Lessons (19 hrs)",
      exams: "2008 - 2016 E.C. Solutions",
      topics: ["Modern Ethiopian State Formation", "World Wars & Global Politics", "Pan-African Movements", "Cold War Dynamics", "Contemporary Ethiopian Developments"],
    },
    {
      title: "Geography",
      icon: Globe2,
      color: "text-[#38bdf8]",
      badgeColor: "bg-[#38bdf8]/15 text-[#38bdf8] border-[#38bdf8]/30",
      units: "6 Core Units",
      videos: "36 Video Lessons (17 hrs)",
      exams: "2008 - 2016 E.C. Solutions",
      topics: ["Map Reading & GIS Fundamentals", "Physical Geography of Ethiopia & Horn", "Human Population & Settlement", "Economic Resources", "Global Climate Change"],
    },
    {
      title: "Economics",
      icon: TrendingUp,
      color: "text-[#8b5cf6]",
      badgeColor: "bg-[#8b5cf6]/15 text-[#8b5cf6] border-[#8b5cf6]/30",
      units: "6 Core Units",
      videos: "38 Video Lessons (18 hrs)",
      exams: "2008 - 2016 E.C. Solutions",
      topics: ["Microeconomic Theory (Supply/Demand)", "Macroeconomic Indicators (GDP, Inflation)", "Fiscal & Monetary Policy", "International Trade", "Ethiopian Economic Structure"],
    },
    {
      title: "English (Advanced)",
      icon: BookText,
      color: "text-[#ec4899]",
      badgeColor: "bg-[#ec4899]/15 text-[#ec4899] border-[#ec4899]/30",
      units: "5 Core Units",
      videos: "30 Video Lessons (14 hrs)",
      exams: "2008 - 2016 E.C. Solutions",
      topics: ["Contextual Vocabulary", "Grammar & Error Identification", "Passage Analysis", "Punctuation & Syntax", "Critical Reasoning"],
    },
    {
      title: "Civics & General Aptitude",
      icon: Landmark,
      color: "text-[#14b8a6]",
      badgeColor: "bg-[#14b8a6]/15 text-[#14b8a6] border-[#14b8a6]/30",
      units: "5 Core Units",
      videos: "26 Video Lessons (12 hrs)",
      exams: "2008 - 2016 E.C. Solutions",
      topics: ["Ethiopian Constitutional Framework", "Democratic Governance", "Logical Fallacies & Syllogisms", "Data Interpretation", "Verbal Aptitude"],
    },
  ];

  const currentCourses = activeStream === "natural" ? naturalCourses : socialCourses;

  return (
    <section id="curriculum" className="py-24 bg-[#050505] relative">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Section Header */}
        <div className="text-center max-w-3xl mx-auto mb-12 space-y-4">
          <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-[#10b981]/10 text-[#10b981] border border-[#10b981]/25 text-xs font-semibold">
            <Layers className="w-3.5 h-3.5" />
            <span>Grade 12 Curriculum Streams</span>
          </div>
          <h2 className="text-3xl sm:text-4xl font-extrabold text-white tracking-tight">
            100% Comprehensive Coverage for Both Streams
          </h2>
          <p className="text-sm sm:text-base text-[#94a3b8]">
            Every single topic mandated by the Ethiopian Ministry of Education, taught by seasoned university prep instructors with step-by-step past national exam walkthroughs.
          </p>

          {/* Stream Switcher Tabs */}
          <div className="inline-flex p-1.5 rounded-2xl bg-[#111115] border border-[#1a1a20] shadow-xl mt-4">
            <button
              onClick={() => setActiveStream("natural")}
              className={`flex items-center gap-2.5 px-6 py-3 rounded-xl text-sm font-bold transition-all ${
                activeStream === "natural"
                  ? "bg-gradient-to-r from-[#10b981] to-[#059669] text-white shadow-lg shadow-[#10b981]/25"
                  : "text-[#94a3b8] hover:text-white"
              }`}
            >
              <Atom className="w-4 h-4" />
              <span>🌿 Natural Science Stream</span>
            </button>
            <button
              onClick={() => setActiveStream("social")}
              className={`flex items-center gap-2.5 px-6 py-3 rounded-xl text-sm font-bold transition-all ${
                activeStream === "social"
                  ? "bg-gradient-to-r from-[#8b5cf6] to-[#6d28d9] text-white shadow-lg shadow-[#8b5cf6]/25"
                  : "text-[#94a3b8] hover:text-white"
              }`}
            >
              <Landmark className="w-4 h-4" />
              <span>🏛️ Social Science Stream</span>
            </button>
          </div>
        </div>

        {/* Subjects Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {currentCourses.map((course, idx) => {
            const Icon = course.icon;
            return (
              <div
                key={idx}
                className="p-6 rounded-2xl bg-[#111115] border border-[#1a1a20] hover:border-[#2c2c35] transition-all hover:translate-y-[-2px] hover:shadow-2xl hover:shadow-black/50 flex flex-col justify-between group"
              >
                <div>
                  {/* Top Row: Icon + Badge */}
                  <div className="flex items-center justify-between mb-4">
                    <div className="w-12 h-12 rounded-xl bg-[#18181e] border border-[#2c2c35] flex items-center justify-center group-hover:scale-110 transition-transform">
                      <Icon className={`w-6 h-6 ${course.color}`} />
                    </div>
                    <span className={`px-2.5 py-1 rounded-full text-xs font-semibold border ${course.badgeColor}`}>
                      {course.units}
                    </span>
                  </div>

                  {/* Course Title */}
                  <h3 className="text-lg font-bold text-white mb-2 group-hover:text-[#10b981] transition-colors">
                    {course.title}
                  </h3>

                  {/* Video & Exam Metadata */}
                  <div className="space-y-1.5 text-xs text-[#94a3b8] mb-4">
                    <div className="flex items-center gap-2">
                      <PlayCircle className="w-3.5 h-3.5 text-[#10b981]" />
                      <span>{course.videos}</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <CheckCircle2 className="w-3.5 h-3.5 text-[#f59e0b]" />
                      <span>Past National Exams: {course.exams}</span>
                    </div>
                    <div className="flex items-center gap-2">
                      <FileText className="w-3.5 h-3.5 text-[#38bdf8]" />
                      <span>Downloadable PDF Chapter Notes Included</span>
                    </div>
                  </div>

                  {/* Topic Pill Tags */}
                  <div className="border-t border-[#1a1a20] pt-3">
                    <span className="text-[11px] font-semibold text-[#71717a] block mb-2 uppercase tracking-wider">
                      Key Topics Covered:
                    </span>
                    <div className="flex flex-wrap gap-1.5">
                      {course.topics.map((topic, tIdx) => (
                        <span
                          key={tIdx}
                          className="px-2 py-0.5 rounded-md bg-[#18181e] text-[#d4d4d8] text-[11px] border border-[#25252e]"
                        >
                          {topic}
                        </span>
                      ))}
                    </div>
                  </div>
                </div>

                {/* Card Footer CTA */}
                <div className="mt-6 pt-4 border-t border-[#1a1a20] flex items-center justify-between text-xs">
                  <span className="text-[#94a3b8]">Offline Downloadable</span>
                  <a
                    href="https://play.google.com/store/apps"
                    target="_blank"
                    rel="noopener noreferrer"
                    className="font-bold text-[#10b981] hover:underline flex items-center gap-1"
                  >
                    <span>Start Learning</span>
                    <span>→</span>
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
