"use client";

import React, { useState } from "react";
import Image from "next/image";
import { 
  Smartphone, 
  BookOpen, 
  Clock, 
  Play, 
  BarChart3, 
  WifiOff, 
  CheckCircle2, 
  ArrowRight,
  Sparkles,
  Download,
  Eye,
  ShieldCheck
} from "lucide-react";

interface AppScreen {
  id: string;
  title: string;
  subtitle: string;
  badge: string;
  icon: React.ElementType;
  image: string;
  highlights: string[];
  description: string;
}

const screens: AppScreen[] = [
  {
    id: "home",
    title: "Course Catalog & Grades",
    subtitle: "Grade 9–12 & Freshman Streams",
    badge: "Personalized Home",
    icon: BookOpen,
    image: "/screenshots/mirkuz_home.jpg",
    highlights: [
      "Natural & Social Science filtered by Grade 9 to 12",
      "Featured interactive master courses with Mirkuz Mascot",
      "One-tap subject selector: Math, Physics, Chemistry, Biology & English",
      "Clean Obsidian Dark UI built for night & daytime study"
    ],
    description:
      "Explore structured Ethiopian curriculum courses categorized by stream and grade. Fast access to video lessons, short summaries, and quizzes."
  },
  {
    id: "exam",
    title: "Simulated Entrance Exam",
    subtitle: "Timed Mock Tests & Questions",
    badge: "Authentic Simulation",
    icon: Clock,
    image: "/screenshots/mirkuz_exam.jpg",
    highlights: [
      "Authentic national entrance countdown timer",
      "60-question interactive jump palette with review flags",
      "High-contrast math & science equation rendering",
      "Instant step-by-step video & text solutions after submission"
    ],
    description:
      "Prepare under real university entrance exam conditions. Take timed past matric exams with authentic pacing, question palettes, and instant feedback."
  },
  {
    id: "video",
    title: "HD Video Lessons & Notes",
    subtitle: "Taught by Senior Ethiopian Examiners",
    badge: "Adaptive Streaming",
    icon: Play,
    image: "/screenshots/mirkuz_video.jpg",
    highlights: [
      "Low-bandwidth adaptive HLS streaming (360p to 720p HD)",
      "Variable playback speeds (1.0x, 1.25x, 1.5x, 2.0x)",
      "Digital blackboard visual derivations & experiments",
      "Downloadable unit summary PDFs for quick revision"
    ],
    description:
      "Clear, conceptual video lessons taught by top Ethiopian high school educators and national examiners. Follow along unit by unit with PDF chapter handouts."
  },
  {
    id: "analytics",
    title: "AI Diagnostics & Score Rank",
    subtitle: "Pinpoint Weak Areas Before Exam Day",
    badge: "Score Prediction",
    icon: BarChart3,
    image: "/screenshots/mirkuz_analytics.jpg",
    highlights: [
      "National percentile ranking & predicted university placement",
      "Chapter-by-chapter mastery progress breakdown",
      "AI weak spot detection with targeted lesson recommendations",
      "Daily study streak tracker & milestone achievements"
    ],
    description:
      "Stop guessing your exam readiness. Mirkuz analytics identifies the exact chapters and formulas pulling your score down so you can fix them early."
  },
  {
    id: "offline",
    title: "100% Offline Study Mode",
    subtitle: "Zero Mobile Data Required",
    badge: "Encrypted Storage",
    icon: WifiOff,
    image: "/screenshots/mirkuz_offline.jpg",
    highlights: [
      "One-click download of full subject chapters & video lessons",
      "Encrypted lightweight storage takes minimal phone space",
      "Take mock exams & watch lessons completely offline",
      "Saves up to 15 GB of mobile data every month"
    ],
    description:
      "Download your study units while on Wi-Fi at home or school, then turn off mobile data and study anywhere in Ethiopia with zero internet connection."
  }
];

export function AppShowcase() {
  const [activeTab, setActiveTab] = useState(0);
  const current = screens[activeTab];

  return (
    <section id="app-tour" className="py-24 bg-[#101730] relative overflow-hidden border-t border-[#1E2B58]">
      {/* Background Ambient Glows */}
      <div className="absolute top-1/3 left-1/4 -translate-x-1/2 w-[600px] h-[600px] bg-[#00B894]/12 rounded-full blur-[160px] pointer-events-none" />
      <div className="absolute bottom-1/4 right-1/4 w-[500px] h-[500px] bg-[#2A3B72]/20 rounded-full blur-[140px] pointer-events-none" />

      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        
        {/* Section Header with Underline */}
        <div className="text-center max-w-3xl mx-auto mb-14 space-y-4">
          <div className="inline-flex items-center gap-2 px-3.5 py-1.5 rounded-full bg-[#00B894]/15 text-[#00B894] border border-[#00B894]/30 text-xs font-bold">
            <Smartphone className="w-3.5 h-3.5" />
            <span>Interactive Mobile Experience</span>
          </div>

          <h2 className="text-3xl sm:text-4xl lg:text-5xl font-extrabold text-white tracking-tight">
            See the Actual{" "}
            <span className="text-[#00B894]">
              Mirkuz Mobile App
            </span>{" "}
            in Action
          </h2>
          <div className="title-underline-center" />

          <p className="text-sm sm:text-base text-[#94A3B8] leading-relaxed pt-2">
            Designed specifically for Ethiopian Grade 9–12 students. Fast, responsive, dark-mode native, and packed with national matric tools.
          </p>
        </div>

        {/* Interactive Feature Selector Tabs */}
        <div className="flex items-center justify-start lg:justify-center gap-2.5 overflow-x-auto pb-4 scrollbar-none mb-12">
          {screens.map((screen, idx) => {
            const Icon = screen.icon;
            const isSelected = activeTab === idx;
            return (
              <button
                key={screen.id}
                onClick={() => setActiveTab(idx)}
                className={`flex items-center gap-2.5 px-4 py-2.5 rounded-2xl text-xs sm:text-sm font-bold whitespace-nowrap transition-all border shrink-0 ${
                  isSelected
                    ? "bg-[#00B894] text-white border-[#00B894] shadow-lg shadow-[#00B894]/25"
                    : "bg-[#162145] text-[#94A3B8] hover:text-white hover:bg-[#1E2B58] border-[#2A3B72]/60"
                }`}
              >
                <Icon className={`w-4 h-4 ${isSelected ? "text-white" : "text-[#00B894]"}`} />
                <span>{screen.title}</span>
              </button>
            );
          })}
        </div>

        {/* Main Showcase Grid */}
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-10 lg:gap-12 items-center">
          
          {/* Left Column: Screen Details & Highlights */}
          <div className="lg:col-span-6 space-y-6 order-2 lg:order-1">
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-[#162145] border border-[#2A3B72] text-xs font-bold text-[#00B894]">
              <Sparkles className="w-3.5 h-3.5" />
              <span>{current.badge}</span>
            </div>

            <h3 className="text-2xl sm:text-3xl font-extrabold text-white tracking-tight">
              {current.title}
            </h3>

            <p className="text-sm sm:text-base text-[#94A3B8] leading-relaxed">
              {current.description}
            </p>

            {/* Highlights List */}
            <div className="space-y-3 pt-2">
              {current.highlights.map((item, i) => (
                <div
                  key={i}
                  className="p-3.5 rounded-2xl bg-[#162145]/90 border border-[#2A3B72]/50 flex items-start gap-3 hover:border-[#00B894]/40 transition-colors"
                >
                  <div className="w-5 h-5 rounded-full bg-[#00B894]/20 flex items-center justify-center text-[#00B894] shrink-0 mt-0.5">
                    <CheckCircle2 className="w-3.5 h-3.5" />
                  </div>
                  <span className="text-xs sm:text-sm text-[#E2E8F0] font-medium leading-normal">
                    {item}
                  </span>
                </div>
              ))}
            </div>

            {/* Direct CTA */}
            <div className="pt-4 flex flex-col sm:flex-row items-stretch sm:items-center gap-3">
              <a
                href="https://play.google.com/store/apps/details?id=et.mirkuz.mobile"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center justify-center gap-2 px-6 py-3.5 rounded-xl bg-[#00B894] hover:bg-[#00A381] text-white font-bold text-xs sm:text-sm shadow-xl shadow-[#00B894]/25 hover:shadow-2xl hover:scale-[1.02] active:scale-[0.98] transition-all"
              >
                <Download className="w-4 h-4" />
                <span>Get Mirkuz on Google Play</span>
              </a>

              <a
                href="https://t.me/mirkuz_exam"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center justify-center gap-2 px-5 py-3.5 rounded-xl bg-[#162145] text-[#CBD5E1] hover:text-white border border-[#2A3B72] font-semibold text-xs sm:text-sm transition-colors"
              >
                <span>Ask on Telegram</span>
                <ArrowRight className="w-4 h-4 text-[#00B894]" />
              </a>
            </div>
          </div>

          {/* Right Column: High-Fidelity Phone Frame with Real Screenshot */}
          <div className="lg:col-span-6 flex justify-center order-1 lg:order-2">
            <div className="relative">
              
              {/* Outer Glow Halo */}
              <div className="absolute -inset-4 bg-gradient-to-tr from-[#00B894]/25 via-[#2A3B72]/20 to-transparent rounded-[52px] blur-2xl opacity-75 pointer-events-none" />

              {/* Realistic Modern Phone Frame */}
              <div className="relative w-[290px] sm:w-[330px] rounded-[48px] p-3 bg-gradient-to-b from-[#2A3B72] via-[#162145] to-[#101730] shadow-2xl shadow-black/90 border border-[#3A4D8C]">
                
                {/* Screen Housing */}
                <div className="relative rounded-[38px] overflow-hidden bg-[#050505] border border-[#181820] shadow-inner aspect-[9/18.5]">
                  
                  {/* Real App Screenshot */}
                  <Image
                    src={current.image}
                    alt={current.title}
                    fill
                    sizes="(max-width: 640px) 290px, 330px"
                    priority
                    className="object-cover transition-opacity duration-300"
                  />

                  {/* Glass Reflection Glare */}
                  <div className="absolute inset-0 bg-gradient-to-tr from-transparent via-white/[0.04] to-transparent pointer-events-none" />
                </div>

                {/* Bottom Home Indicator Bar */}
                <div className="w-28 h-1 bg-white/30 rounded-full mx-auto mt-2" />
              </div>

              {/* Floating Live Feature Tag */}
              <div className="absolute -bottom-4 -left-4 sm:-left-6 bg-[#162145]/95 backdrop-blur-md p-3 rounded-2xl border border-[#2A3B72] shadow-2xl flex items-center gap-3">
                <div className="w-8 h-8 rounded-xl bg-[#00B894]/20 flex items-center justify-center text-[#00B894]">
                  <ShieldCheck className="w-4 h-4" />
                </div>
                <div>
                  <div className="text-[11px] font-bold text-white">Ethiopian Curriculum Aligned</div>
                  <div className="text-[9.5px] text-[#94A3B8]">Updated for 2017/2018 E.C.</div>
                </div>
              </div>

            </div>
          </div>

        </div>

      </div>
    </section>
  );
}
