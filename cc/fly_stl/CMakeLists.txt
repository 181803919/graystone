cmake_minimum_required(VERSION 2.8)

project(fly_stl)
set(CMAKE_INCLUDE_CURRENT_DIR ON)

option (USE_MYSKIP
	   "Use provided skip implementation" ON)

configure_file (
  "${PROJECT_SOURCE_DIR}/config.h.in"
  "${PROJECT_SOURCE_DIR}/config.h"
  )

if (USE_MYSKIP)
  include_directories ("${PROJECT_SOURCE_DIR}/skip")
  add_subdirectory (skip)
  set (EXTRA_LIBS ${EXTRA_LIBS} SkipLib)
endif (USE_MYSKIP)

aux_source_directory(. fly_stl_src)

add_executable(fly_stl_test ${fly_stl_src})
target_link_libraries(fly_stl_test ${EXTRA_LIBS})
#set(CMAKE_BUILD_TYPE "Debug")
#set(CMAKE_CXX_FLAGS_DEBUG "$ENV{CXXFLAGS} -O0 -Wall -g -ggdb")
#set(CMAKE_CXX_FLAGS_RELEASE "$ENV{CXXFLAGS} -O3 -Wall")
