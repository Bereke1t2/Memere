"use client";

import React, { useState } from "react";
import Image from "next/image";
import { 
  ArrowLeft, 
  ArrowRight, 
  Star, 
  GraduationCap, 
  Smartphone,
  ShieldCheck,
  Send
} from "lucide-react";

export function TestimonialSpotlight() {
  const [activeIdx, setActiveIdx] = useState(0);

  const testimonials = [
    {
      title: "Scored 642 and Placed into AAU Medicine!",
      quote:
        "There is only one word I can use to describe my experience with Mirkuz: partner. In the truest sense of a high-performing partnership. The past exam video derivations uncovered shortcuts my regular classes never taught. Best of all, having all video masterclasses available offline meant I could study late at night without worrying about weak mobile internet or data costs.",
      name: "Kalkidan Bekele",
      role: "AAU Medical School Candidate • 642/700 (2016 E.C.)",
      country: "ET",
      flag: "🇪🇹",
      avatarBg: "from-[#00B894] to-[#00A381]",
    },
    {
      title: "Top 1% Score in Mathematics & Physics",
      quote:
        "The timed mock exam simulator gave me the exact timing discipline I needed for the real matric. The 60-question jump palette and immediate step-by-step video solutions transformed how I solved calculus limits and rotational mechanics questions.",
      name: "Yohannes Tadesse",
      role: "AASTU Software Engineering • 618/700",
      country: "ET",
      flag: "🇪🇹",
      avatarBg: "from-[#1E2B58] to-[#2A3B72]",
    },
    {
      title: "Zero Internet Needed in Dormitory Study",
      quote:
        "Living in a high school boarding dorm with spotty mobile data, Mirkuz's offline download feature was a lifesaver. I downloaded all Grade 12 Chemistry and Biology units and studied every evening completely offline.",
      name: "Selamawit Girma",
      role: "Jimma University Health Science • 594/700",
      country: "ET",
      flag: "🇪🇹",
      avatarBg: "from-[#00B894] to-[#1E2B58]",
    },
    {
      title: "Social Science & Aptitude Mastery",
      quote:
        "The aptitude speed tactics alone saved me at least 25 minutes on exam day. Economics models and Ethiopian history summaries were crisp, high-yield, and perfectly matched the 2017 E.C. syllabus.",
      name: "Abel Tesfaye",
      role: "Hawassa University Economics • 582/700",
      country: "ET",
      flag: "🇪🇹",
      avatarBg: "from-[#2A3B72] to-[#00B894]",
    },
  ];

  const handlePrev = () => {
    setActiveIdx((prev) => (prev === 0 ? testimonials.length - 1 : prev - 1));
  };

  const handleNext = () => {
    setActiveIdx((prev) => (prev === testimonials.length - 1 ? 0 : prev + 1));
  };

  return (
    <section id="testimonials" className="bg-white py-20 md:py-28 border-b border-[#F1F5F9] relative overflow-hidden">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        
        {/* Header with Circular Navigation Arrows Matching Frame 008 */}
        <div className="flex flex-col sm:flex-row sm:items-end justify-between mb-14 gap-6">
          <div>
            <h2 className="text-3xl sm:text-4xl lg:text-5xl font-extrabold text-[#1E2B58] tracking-tight">
              What students say
            </h2>
            <div className="title-underline" />
            <p className="text-xs sm:text-sm text-[#64748B] mt-4">
              See what top scorers say about preparing with Mirkuz
            </p>
          </div>

          {/* Circular Arrow Buttons (Frame 008) */}
          <div className="flex items-center gap-3 self-start sm:self-auto">
            <button
              onClick={handlePrev}
              aria-label="Previous testimonial"
              className="w-12 h-12 rounded-full border border-[#E2E8F0] bg-white hover:bg-[#F8FAFC] flex items-center justify-center text-[#1E2B58] hover:text-[#00B894] hover:border-[#00B894] shadow-sm transition-all"
            >
              <ArrowLeft className="w-4 h-4" />
            </button>
            <button
              onClick={handleNext}
              aria-label="Next testimonial"
              className="w-12 h-12 rounded-full border border-[#E2E8F0] bg-white hover:bg-[#F8FAFC] flex items-center justify-center text-[#1E2B58] hover:text-[#00B894] hover:border-[#00B894] shadow-sm transition-all"
            >
              <ArrowRight className="w-4 h-4" />
            </button>
          </div>
        </div>

        {/* Testimonial Cards Carousel Grid Matching Frame 008 */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-7 mb-12">
          {testimonials.slice(0, 3).map((item, idx) => (
            <div
              key={idx}
              className="bg-white rounded-3xl p-8 border border-[#E2E8F0] shadow-dribbble hover:shadow-dribbble-lg transition-all duration-300 hover:-translate-y-1 flex flex-col justify-between"
            >
              <div>
                {/* 5 Gold Stars */}
                <div className="flex items-center gap-1 mb-4 text-[#F59E0B]">
                  {[...Array(5)].map((_, i) => (
                    <Star key={i} className="w-4 h-4 fill-current" />
                  ))}
                </div>

                {/* Punchy Title */}
                <h3 className="text-base sm:text-lg font-extrabold text-[#1E2B58] mb-3 leading-snug">
                  {item.title}
                </h3>

                {/* Italic Quote */}
                <p className="text-xs sm:text-sm text-[#64748B] italic leading-relaxed">
                  “{item.quote}”
                </p>
              </div>

              {/* Author Row with Avatar and Country Flag Badge */}
              <div className="flex items-center justify-between pt-6 mt-6 border-t border-[#F1F5F9]">
                <div className="flex items-center gap-3">
                  <div className={`w-11 h-11 rounded-full bg-gradient-to-br ${item.avatarBg} text-white flex items-center justify-center font-bold text-sm shadow-sm`}>
                    {item.name[0]}
                  </div>
                  <div>
                    <h4 className="text-sm font-bold text-[#1E2B58]">{item.name}</h4>
                    <p className="text-[11px] text-[#94A3B8] max-w-[170px] truncate">{item.role}</p>
                  </div>
                </div>

                {/* Ethiopian Flag Badge (🇪🇹) */}
                <span className="text-xl select-none" title="Ethiopia">
                  {item.flag}
                </span>
              </div>
            </div>
          ))}
        </div>

        {/* Carousel Pagination Dots Matching Frame 008 */}
        <div className="flex items-center justify-center gap-2 mb-20">
          {testimonials.map((_, i) => (
            <button
              key={i}
              onClick={() => setActiveIdx(i)}
              aria-label={`Go to slide ${i + 1}`}
              className={`rounded-full transition-all ${
                activeIdx === i
                  ? "w-8 h-2.5 bg-[#00B894]"
                  : "w-2.5 h-2.5 bg-[#CBD5E1] hover:bg-[#94A3B8]"
              }`}
            />
          ))}
        </div>

        {/* Bottom Navy Callout Banner Matching Frame 007 & 008 */}
        <div className="rounded-3xl bg-[#1E2B58] p-10 sm:p-16 text-center text-white shadow-2xl shadow-[#1E2B58]/25 relative overflow-hidden">
          {/* Subtle Ambient Radial Glow */}
          <div className="absolute top-0 right-0 w-80 h-80 bg-[#00B894]/20 rounded-full blur-3xl pointer-events-none" />
          <div className="absolute bottom-0 left-0 w-80 h-80 bg-[#2A3B72]/40 rounded-full blur-2xl pointer-events-none" />

          <div className="relative z-10 max-w-3xl mx-auto space-y-6">
            <h3 className="text-2xl sm:text-4xl lg:text-5xl font-extrabold tracking-tight leading-tight text-white">
              Start acting based on real data. <br />
              Different preparation, different impact.
            </h3>

            <p className="text-sm sm:text-base text-[#CBD5E1] max-w-xl mx-auto">
              Download Mirkuz on Android today. Study past entrance exams, watch examiner masterclasses, and take timed mock tests offline.
            </p>

            <div className="pt-2">
              <a
                href="https://play.google.com/store/apps/details?id=et.mirkuz.mobile"
                target="_blank"
                rel="noopener noreferrer"
                className="inline-flex items-center justify-center gap-2 px-8 py-4 rounded-xl bg-[#00B894] hover:bg-[#00A381] text-white font-bold text-sm sm:text-base shadow-xl shadow-[#00B894]/30 hover:scale-[1.02] active:scale-[0.98] transition-all"
              >
                <Smartphone className="w-4 h-4" />
                <span>Get Started</span>
              </a>
            </div>
          </div>
        </div>

      </div>
    </section>
  );
}
