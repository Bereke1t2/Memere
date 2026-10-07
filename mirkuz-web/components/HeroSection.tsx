"use client";

import React, { useState, useEffect } from "react";
import Image from "next/image";
import { 
  ArrowRight, 
  CheckCircle2, 
  Sparkles, 
  TrendingUp, 
  Download, 
  Play, 
  BookOpen,
  Award,
  Smartphone,
  ShieldCheck
} from "lucide-react";

export function HeroSection() {
  const [activeSlide, setActiveSlide] = useState(0);

  const slides = [
    {
      titlePrefix: "Master Ethiopia's Grade 12",
      titleHighlight: "National Entrance Exam",
      subtitle:
        "Where top scores begin: excel in your matric exam with senior examiner video lessons, 10,000+ past entrance questions, and 100% data-free offline study.",
      ctaPrimary: "Download App",
      ctaSecondary: "Explore Curriculum",
      badge: "2017 E.C. Ministry Curriculum Aligned",
      accent: "#4CAF4F"
    },
    {
      titlePrefix: "Master Your National Exam",
      titleHighlight: "with Mirkuz Prep",
      subtitle:
        "10,000+ national entrance questions (2008–2016 E.C.) with step-by-step video solutions, timed exam simulator, and instant score prediction.",
      ctaPrimary: "Download App",
      ctaSecondary: "Try Mock Exam",
      badge: "Over 50,000 Active Students",
      accent: "#4CAF4F"
    },
    {
      titlePrefix: "100% Offline Study Mode",
      titleHighlight: "Zero Mobile Data",
      subtitle:
        "Save all chapters and video courses directly to your phone. Turn off mobile data and study anywhere across Ethiopia without spending Birr on internet packages.",
      ctaPrimary: "Download App",
      ctaSecondary: "See How It Works",
      badge: "Built for Low Bandwidth",
      accent: "#4CAF4F"
    },
  ];

  // Auto rotate slides every 7 seconds
  useEffect(() => {
    const timer = setInterval(() => {
      setActiveSlide((prev) => (prev + 1) % slides.length);
    }, 7000);
    return () => clearInterval(timer);
  }, [slides.length]);

  const current = slides[activeSlide];

  return (
    <section id="hero" className="bg-[#F5F7FA] pt-32 pb-20 md:pt-40 md:pb-28 relative overflow-hidden transition-colors duration-500">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        <div className="grid grid-cols-1 lg:grid-cols-12 gap-12 lg:gap-8 items-center min-h-[460px]">
          
          {/* Left Column: Typography & CTAs */}
          <div className="lg:col-span-7 flex flex-col items-start space-y-6">
            
            {/* Top Minimal Badge */}
            <div className="inline-flex items-center gap-2 px-3 py-1 rounded-full bg-white border border-[#E8ECF2] shadow-xs text-xs font-semibold text-[#4D4D4D]">
              <span className="w-2 h-2 rounded-full bg-[#4CAF4F] animate-pulse"></span>
              <span>{current.badge}</span>
            </div>

            {/* Main Headline */}
            <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold tracking-tight text-[#4D4D4D] leading-[1.15]">
              {current.titlePrefix} <br />
              <span className="text-[#4CAF4F]">{current.titleHighlight}</span>
            </h1>

            {/* Subtitle */}
            <p className="text-base sm:text-lg text-[#717171] max-w-xl leading-relaxed">
              {current.subtitle}
            </p>

            {/* Action Buttons */}
            <div className="flex flex-wrap items-center gap-4 pt-2">
              <a
                href="https://play.google.com/store/apps/details?id=et.mirkuz.mobile"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center justify-center gap-2 px-8 py-3.5 rounded-md bg-[#4CAF4F] hover:bg-[#388E3C] text-white font-medium text-base shadow-sm hover:shadow-md hover:-translate-y-0.5 active:translate-y-0 transition-all"
              >
                <span>{current.ctaPrimary}</span>
                <ArrowRight className="w-4 h-4" />
              </a>

              <a
                href="#curriculum"
                className="inline-flex items-center justify-center gap-2 px-7 py-3.5 rounded-md bg-white hover:bg-neutral-50 text-[#4D4D4D] border border-[#D5E0D5] font-medium text-base shadow-2xs hover:border-[#4CAF4F] transition-all"
              >
                <span>{current.ctaSecondary}</span>
              </a>
            </div>

            {/* Trust points */}
            <div className="flex flex-wrap items-center gap-6 pt-4 text-xs text-[#717171]">
              <div className="flex items-center gap-1.5">
                <CheckCircle2 className="w-4 h-4 text-[#4CAF4F]" />
                <span>Free Starter Access</span>
              </div>
              <div className="flex items-center gap-1.5">
                <CheckCircle2 className="w-4 h-4 text-[#4CAF4F]" />
                <span>Natural & Social Streams</span>
              </div>
              <div className="flex items-center gap-1.5">
                <CheckCircle2 className="w-4 h-4 text-[#4CAF4F]" />
                <span>Telebirr & Chapa Ready</span>
              </div>
            </div>

          </div>

          {/* Right Column: Clean Vector Hero Illustration */}
          <div className="lg:col-span-5 relative flex items-center justify-center">
            
            {/* Minimal Vector Illustration Canvas */}
            <div className="relative w-full max-w-[440px] aspect-4/3 flex items-center justify-center">
              
              {/* Modern Vector Computer & Study Setup */}
              <svg viewBox="0 0 500 400" className="w-full h-auto drop-shadow-md">
                <defs>
                  <linearGradient id="screenGrad" x1="0" y1="0" x2="1" y2="1">
                    <stop offset="0%" stopColor="#4CAF4F" stopOpacity="0.15" />
                    <stop offset="100%" stopColor="#E8F5E9" stopOpacity="0.8" />
                  </linearGradient>
                  <linearGradient id="barGrad" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="0%" stopColor="#4CAF4F" />
                    <stop offset="100%" stopColor="#388E3C" />
                  </linearGradient>
                </defs>

                {/* Desk Surface Shadow */}
                <ellipse cx="250" cy="365" rx="190" ry="14" fill="#E8ECF2" />

                {/* Monitor Stand */}
                <rect x="235" y="270" width="30" height="60" rx="4" fill="#89939E" />
                <path d="M190 330 L310 330 L320 345 L180 345 Z" fill="#717171" />

                {/* Monitor Frame */}
                <rect x="90" y="70" width="320" height="210" rx="12" fill="#263238" />
                {/* Screen Bezel / Display */}
                <rect x="100" y="80" width="300" height="185" rx="6" fill="url(#screenGrad)" />

                {/* Screen Top Bar */}
                <rect x="100" y="80" width="300" height="24" rx="4" fill="#FFFFFF" />
                <circle cx="115" cy="92" r="4" fill="#E53835" />
                <circle cx="127" cy="92" r="4" fill="#FBC02D" />
                <circle cx="139" cy="92" r="4" fill="#4CAF4F" />
                <rect x="155" y="87" width="120" height="10" rx="5" fill="#F5F7FA" />

                {/* Screen Content: Charts & Analytics */}
                {/* Left Mini Sidebar */}
                <rect x="110" y="114" width="50" height="140" rx="4" fill="#FFFFFF" />
                <rect x="116" y="122" width="38" height="6" rx="3" fill="#4CAF4F" />
                <rect x="116" y="134" width="30" height="4" rx="2" fill="#ABBED1" />
                <rect x="116" y="144" width="34" height="4" rx="2" fill="#ABBED1" />
                <rect x="116" y="154" width="28" height="4" rx="2" fill="#ABBED1" />
                <rect x="116" y="164" width="32" height="4" rx="2" fill="#ABBED1" />

                {/* Main Graph Area */}
                <rect x="170" y="114" width="220" height="85" rx="6" fill="#FFFFFF" />
                <text x="180" y="132" fill="#4D4D4D" fontSize="10" fontWeight="bold">Exam Score Progression</text>
                {/* Bar chart lines */}
                <line x1="180" y1="185" x2="375" y2="185" stroke="#E8ECF2" strokeWidth="1" />
                <rect x="195" y="155" width="18" height="30" rx="2" fill="#ABBED1" />
                <rect x="225" y="145" width="18" height="40" rx="2" fill="#ABBED1" />
                <rect x="255" y="135" width="18" height="50" rx="2" fill="#4CAF4F" />
                <rect x="285" y="125" width="18" height="60" rx="2" fill="url(#barGrad)" />
                <rect x="315" y="118" width="18" height="67" rx="2" fill="#4CAF4F" />
                <rect x="345" y="112" width="18" height="73" rx="2" fill="#2E7D32" />

                {/* Screen Cards Bottom Row */}
                <rect x="170" y="208" width="105" height="46" rx="4" fill="#FFFFFF" />
                <circle cx="185" cy="225" r="8" fill="#E8F5E9" />
                <path d="M182 225 L184 227 L189 222" stroke="#4CAF4F" strokeWidth="2" fill="none" />
                <rect x="200" y="220" width="60" height="5" rx="2.5" fill="#4D4D4D" />
                <rect x="200" y="228" width="40" height="4" rx="2" fill="#89939E" />

                <rect x="285" y="208" width="105" height="46" rx="4" fill="#FFFFFF" />
                <circle cx="300" cy="225" r="8" fill="#E8F5E9" />
                <text x="296" y="229" fill="#4CAF4F" fontSize="10" fontWeight="bold">★</text>
                <rect x="315" y="220" width="60" height="5" rx="2.5" fill="#4D4D4D" />
                <rect x="315" y="228" width="45" height="4" rx="2" fill="#89939E" />

                {/* Student Character on the Right */}
                {/* Body / Torso */}
                <path d="M380 260 C380 230, 420 230, 420 260 L425 350 L375 350 Z" fill="#4CAF4F" />
                {/* Arms */}
                <path d="M380 250 L350 220 L360 215 L390 240 Z" fill="#388E3C" />
                <circle cx="348" cy="216" r="6" fill="#F8CBA6" />
                {/* Head */}
                <circle cx="400" cy="210" r="18" fill="#F8CBA6" />
                {/* Hair */}
                <path d="M384 205 C384 190, 416 190, 416 205 C416 195, 384 195, 384 205 Z" fill="#263238" />
                {/* Glasses / Face detail */}
                <rect x="390" y="208" width="9" height="5" rx="1" fill="#263238" />
                <rect x="403" y="208" width="9" height="5" rx="1" fill="#263238" />
                {/* Pants */}
                <rect x="382" y="340" width="16" height="40" fill="#263238" />
                <rect x="402" y="340" width="16" height="40" fill="#263238" />
                {/* Shoes */}
                <rect x="375" y="375" width="23" height="8" rx="4" fill="#717171" />
                <rect x="402" y="375" width="23" height="8" rx="4" fill="#717171" />

                {/* Floating Elements / Success Indicator */}
                <g className="animate-bounce" style={{ animationDuration: "3s" }}>
                  <rect x="40" y="140" width="90" height="34" rx="6" fill="#FFFFFF" filter="drop-shadow(0 4px 6px rgba(0,0,0,0.08))" />
                  <circle cx="56" cy="157" r="7" fill="#E8F5E9" />
                  <path d="M53 157 L55 159 L60 154" stroke="#4CAF4F" strokeWidth="1.5" fill="none" />
                  <text x="68" y="154" fill="#4D4D4D" fontSize="9" fontWeight="bold">+18.5% Score</text>
                  <text x="68" y="164" fill="#717171" fontSize="7">AAU Medical Cutoff</text>
                </g>
              </svg>

            </div>

          </div>

        </div>

        {/* Carousel Pagination Dots */}
        <div className="flex items-center justify-center gap-2 pt-10">
          {slides.map((_, idx) => (
            <button
              key={idx}
              onClick={() => setActiveSlide(idx)}
              className={`transition-all rounded-full ${
                activeSlide === idx
                  ? "w-8 h-2.5 bg-[#4CAF4F]"
                  : "w-2.5 h-2.5 bg-[#ABBED1] hover:bg-[#717171]"
              }`}
              aria-label={`Go to slide ${idx + 1}`}
            />
          ))}
        </div>

      </div>
    </section>
  );
}
