/*
 Navicat Premium Data Transfer

 Source Server         : [uat] ciliumlogs
 Source Server Type    : PostgreSQL
 Source Server Version : 130006
 Source Host           : localhost:5432
 Source Catalog        : ciliumlogs
 Source Schema         : public

 Target Server Type    : PostgreSQL
 Target Server Version : 130006
 File Encoding         : 65001

 Date: 04/06/2022 22:47:32
*/


-- ----------------------------
-- Table structure for flows
-- ----------------------------
DROP TABLE IF EXISTS "public"."flows";
CREATE TABLE "public"."flows" (
  "id" int4 NOT NULL DEFAULT nextval('flows_id_seq'::regclass),
  "verdict" varchar COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "ip_source" varchar COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "ip_destination" varchar COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "l4_protocol" varchar COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "l4_source_port" int4 DEFAULT 0,
  "l4_destination_port" int4 DEFAULT 0,
  "source_pod_name" varchar COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "source_namespace" varchar COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "destination_pod_name" varchar COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "destination_namespace" varchar COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "layer" varchar COLLATE "pg_catalog"."default" DEFAULT ''::character varying,
  "checksum" varchar COLLATE "pg_catalog"."default" NOT NULL DEFAULT ''::character varying,
  "flow_count" int4 DEFAULT 1,
  "metadata" json NOT NULL,
  "cluster_name" varchar(128) COLLATE "pg_catalog"."default" DEFAULT ''::character varying
)
;
ALTER TABLE "public"."flows" OWNER TO "cilium";

-- ----------------------------
-- Primary Key structure for table flows
-- ----------------------------
ALTER TABLE "public"."flows" ADD CONSTRAINT "flows_pkey" PRIMARY KEY ("id");
