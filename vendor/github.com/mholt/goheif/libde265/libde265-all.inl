#if _WIN32
#include "extra/win32cond.c"
#define HAVE___MINGW_ALIGNED_MALLOC 1
#else
#define HAVE_POSIX_MEMALIGN 1
#endif

#include "alloc_pool.cc"
#include "bitstream.cc"
#include "cabac.cc"
#include "configparam.cc"
#include "contextmodel.cc"
#include "de265.cc"
#include "deblock.cc"
#include "decctx.cc"
#include "dpb.cc"
// #include "en265.cc"
#include "fallback-dct.cc"
#include "fallback-motion.cc"
#include "fallback.cc"
#include "image-io.cc"
#include "image.cc"
#include "intrapred.cc"
#include "md5.cc"
#include "motion.cc"
#include "nal-parser.cc"
#include "nal.cc"
#include "pps.cc"
#include "quality.cc"
#include "refpic.cc"
#include "sao.cc"
#include "scan.cc"
#include "sei.cc"
#include "slice.cc"
#include "sps.cc"
#include "threads.cc"
#include "transform.cc"
#include "util.cc"
#include "visualize.cc"
#include "vps.cc"
#include "vui.cc"

#ifdef HAVE_SSE4_1
#include "x86/sse-dct.cc"
#include "x86/sse-motion.cc"
#include "x86/sse.cc"
#endif

#ifdef HAVE_ARM
#include "arm/arm.cc"
#endif


