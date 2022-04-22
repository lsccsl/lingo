// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: google/protobuf/unittest_import_public_lite.proto

#ifndef GOOGLE_PROTOBUF_INCLUDED_google_2fprotobuf_2funittest_5fimport_5fpublic_5flite_2eproto
#define GOOGLE_PROTOBUF_INCLUDED_google_2fprotobuf_2funittest_5fimport_5fpublic_5flite_2eproto

#include <limits>
#include <string>

#include <google/protobuf/port_def.inc>
#if PROTOBUF_VERSION < 3019000
#error This file was generated by a newer version of protoc which is
#error incompatible with your Protocol Buffer headers. Please update
#error your headers.
#endif
#if 3019002 < PROTOBUF_MIN_PROTOC_VERSION
#error This file was generated by an older version of protoc which is
#error incompatible with your Protocol Buffer headers. Please
#error regenerate this file with a newer version of protoc.
#endif

#include <google/protobuf/port_undef.inc>
#include <google/protobuf/io/coded_stream.h>
#include <google/protobuf/arena.h>
#include <google/protobuf/arenastring.h>
#include <google/protobuf/generated_message_table_driven.h>
#include <google/protobuf/generated_message_util.h>
#include <google/protobuf/metadata_lite.h>
#include <google/protobuf/message_lite.h>
#include <google/protobuf/repeated_field.h>  // IWYU pragma: export
#include <google/protobuf/extension_set.h>  // IWYU pragma: export
// @@protoc_insertion_point(includes)
#include <google/protobuf/port_def.inc>
#define PROTOBUF_INTERNAL_EXPORT_google_2fprotobuf_2funittest_5fimport_5fpublic_5flite_2eproto
PROTOBUF_NAMESPACE_OPEN
namespace internal {
class AnyMetadata;
}  // namespace internal
PROTOBUF_NAMESPACE_CLOSE

// Internal implementation detail -- do not use these members.
struct TableStruct_google_2fprotobuf_2funittest_5fimport_5fpublic_5flite_2eproto {
  static const ::PROTOBUF_NAMESPACE_ID::internal::ParseTableField entries[]
    PROTOBUF_SECTION_VARIABLE(protodesc_cold);
  static const ::PROTOBUF_NAMESPACE_ID::internal::AuxiliaryParseTableField aux[]
    PROTOBUF_SECTION_VARIABLE(protodesc_cold);
  static const ::PROTOBUF_NAMESPACE_ID::internal::ParseTable schema[1]
    PROTOBUF_SECTION_VARIABLE(protodesc_cold);
  static const ::PROTOBUF_NAMESPACE_ID::internal::FieldMetadata field_metadata[];
  static const ::PROTOBUF_NAMESPACE_ID::internal::SerializationTable serialization_table[];
  static const uint32_t offsets[];
};
namespace protobuf_unittest_import {
class PublicImportMessageLite;
struct PublicImportMessageLiteDefaultTypeInternal;
extern PublicImportMessageLiteDefaultTypeInternal _PublicImportMessageLite_default_instance_;
}  // namespace protobuf_unittest_import
PROTOBUF_NAMESPACE_OPEN
template<> ::protobuf_unittest_import::PublicImportMessageLite* Arena::CreateMaybeMessage<::protobuf_unittest_import::PublicImportMessageLite>(Arena*);
PROTOBUF_NAMESPACE_CLOSE
namespace protobuf_unittest_import {

// ===================================================================

class PublicImportMessageLite final :
    public ::PROTOBUF_NAMESPACE_ID::MessageLite /* @@protoc_insertion_point(class_definition:protobuf_unittest_import.PublicImportMessageLite) */ {
 public:
  inline PublicImportMessageLite() : PublicImportMessageLite(nullptr) {}
  ~PublicImportMessageLite() override;
  explicit constexpr PublicImportMessageLite(::PROTOBUF_NAMESPACE_ID::internal::ConstantInitialized);

  PublicImportMessageLite(const PublicImportMessageLite& from);
  PublicImportMessageLite(PublicImportMessageLite&& from) noexcept
    : PublicImportMessageLite() {
    *this = ::std::move(from);
  }

  inline PublicImportMessageLite& operator=(const PublicImportMessageLite& from) {
    CopyFrom(from);
    return *this;
  }
  inline PublicImportMessageLite& operator=(PublicImportMessageLite&& from) noexcept {
    if (this == &from) return *this;
    if (GetOwningArena() == from.GetOwningArena()
  #ifdef PROTOBUF_FORCE_COPY_IN_MOVE
        && GetOwningArena() != nullptr
  #endif  // !PROTOBUF_FORCE_COPY_IN_MOVE
    ) {
      InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }

  inline const std::string& unknown_fields() const {
    return _internal_metadata_.unknown_fields<std::string>(::PROTOBUF_NAMESPACE_ID::internal::GetEmptyString);
  }
  inline std::string* mutable_unknown_fields() {
    return _internal_metadata_.mutable_unknown_fields<std::string>();
  }

  static const PublicImportMessageLite& default_instance() {
    return *internal_default_instance();
  }
  static inline const PublicImportMessageLite* internal_default_instance() {
    return reinterpret_cast<const PublicImportMessageLite*>(
               &_PublicImportMessageLite_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    0;

  friend void swap(PublicImportMessageLite& a, PublicImportMessageLite& b) {
    a.Swap(&b);
  }
  inline void Swap(PublicImportMessageLite* other) {
    if (other == this) return;
  #ifdef PROTOBUF_FORCE_COPY_IN_SWAP
    if (GetOwningArena() != nullptr &&
        GetOwningArena() == other->GetOwningArena()) {
   #else  // PROTOBUF_FORCE_COPY_IN_SWAP
    if (GetOwningArena() == other->GetOwningArena()) {
  #endif  // !PROTOBUF_FORCE_COPY_IN_SWAP
      InternalSwap(other);
    } else {
      ::PROTOBUF_NAMESPACE_ID::internal::GenericSwap(this, other);
    }
  }
  void UnsafeArenaSwap(PublicImportMessageLite* other) {
    if (other == this) return;
    GOOGLE_DCHECK(GetOwningArena() == other->GetOwningArena());
    InternalSwap(other);
  }

  // implements Message ----------------------------------------------

  PublicImportMessageLite* New(::PROTOBUF_NAMESPACE_ID::Arena* arena = nullptr) const final {
    return CreateMaybeMessage<PublicImportMessageLite>(arena);
  }
  void CheckTypeAndMergeFrom(const ::PROTOBUF_NAMESPACE_ID::MessageLite& from)  final;
  void CopyFrom(const PublicImportMessageLite& from);
  void MergeFrom(const PublicImportMessageLite& from);
  PROTOBUF_ATTRIBUTE_REINITIALIZES void Clear() final;
  bool IsInitialized() const final;

  size_t ByteSizeLong() const final;
  const char* _InternalParse(const char* ptr, ::PROTOBUF_NAMESPACE_ID::internal::ParseContext* ctx) final;
  uint8_t* _InternalSerialize(
      uint8_t* target, ::PROTOBUF_NAMESPACE_ID::io::EpsCopyOutputStream* stream) const final;
  int GetCachedSize() const final { return _cached_size_.Get(); }

  private:
  void SharedCtor();
  void SharedDtor();
  void SetCachedSize(int size) const;
  void InternalSwap(PublicImportMessageLite* other);

  private:
  friend class ::PROTOBUF_NAMESPACE_ID::internal::AnyMetadata;
  static ::PROTOBUF_NAMESPACE_ID::StringPiece FullMessageName() {
    return "protobuf_unittest_import.PublicImportMessageLite";
  }
  protected:
  explicit PublicImportMessageLite(::PROTOBUF_NAMESPACE_ID::Arena* arena,
                       bool is_message_owned = false);
  private:
  static void ArenaDtor(void* object);
  inline void RegisterArenaDtor(::PROTOBUF_NAMESPACE_ID::Arena* arena);
  public:

  std::string GetTypeName() const final;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  enum : int {
    kEFieldNumber = 1,
  };
  // optional int32 e = 1;
  bool has_e() const;
  private:
  bool _internal_has_e() const;
  public:
  void clear_e();
  int32_t e() const;
  void set_e(int32_t value);
  private:
  int32_t _internal_e() const;
  void _internal_set_e(int32_t value);
  public:

  // @@protoc_insertion_point(class_scope:protobuf_unittest_import.PublicImportMessageLite)
 private:
  class _Internal;

  template <typename T> friend class ::PROTOBUF_NAMESPACE_ID::Arena::InternalHelper;
  typedef void InternalArenaConstructable_;
  typedef void DestructorSkippable_;
  ::PROTOBUF_NAMESPACE_ID::internal::HasBits<1> _has_bits_;
  mutable ::PROTOBUF_NAMESPACE_ID::internal::CachedSize _cached_size_;
  int32_t e_;
  friend struct ::TableStruct_google_2fprotobuf_2funittest_5fimport_5fpublic_5flite_2eproto;
};
// ===================================================================


// ===================================================================

#ifdef __GNUC__
  #pragma GCC diagnostic push
  #pragma GCC diagnostic ignored "-Wstrict-aliasing"
#endif  // __GNUC__
// PublicImportMessageLite

// optional int32 e = 1;
inline bool PublicImportMessageLite::_internal_has_e() const {
  bool value = (_has_bits_[0] & 0x00000001u) != 0;
  return value;
}
inline bool PublicImportMessageLite::has_e() const {
  return _internal_has_e();
}
inline void PublicImportMessageLite::clear_e() {
  e_ = 0;
  _has_bits_[0] &= ~0x00000001u;
}
inline int32_t PublicImportMessageLite::_internal_e() const {
  return e_;
}
inline int32_t PublicImportMessageLite::e() const {
  // @@protoc_insertion_point(field_get:protobuf_unittest_import.PublicImportMessageLite.e)
  return _internal_e();
}
inline void PublicImportMessageLite::_internal_set_e(int32_t value) {
  _has_bits_[0] |= 0x00000001u;
  e_ = value;
}
inline void PublicImportMessageLite::set_e(int32_t value) {
  _internal_set_e(value);
  // @@protoc_insertion_point(field_set:protobuf_unittest_import.PublicImportMessageLite.e)
}

#ifdef __GNUC__
  #pragma GCC diagnostic pop
#endif  // __GNUC__

// @@protoc_insertion_point(namespace_scope)

}  // namespace protobuf_unittest_import

// @@protoc_insertion_point(global_scope)

#include <google/protobuf/port_undef.inc>
#endif  // GOOGLE_PROTOBUF_INCLUDED_GOOGLE_PROTOBUF_INCLUDED_google_2fprotobuf_2funittest_5fimport_5fpublic_5flite_2eproto
