import { useMemo, useState } from "react";
import { useNavigate } from "react-router-dom";
import {
  Button,
  ConfirmModal,
  DetailList,
  Drawer,
  EmptyState,
  Field,
  Input,
  Modal,
  PageHeader,
  SurfaceCard,
  StatusPill,
} from "../components/ui";
import { useAppData } from "../contexts/AppDataContext";
import type { BanRecord } from "../types";

export function BansPage() {
  const navigate = useNavigate();
  const { banRecords, createBanRecord, deleteBanRecord } = useAppData();
  const [selected, setSelected] = useState<BanRecord | null>(null);
  const [open, setOpen] = useState(false);
  const [keyword, setKeyword] = useState("");
  const [pendingDelete, setPendingDelete] = useState<BanRecord | null>(null);
  const [namespace, setNamespace] = useState("");
  const [reason, setReason] = useState("");
  const [banStartTime, setBanStartTime] = useState("");
  const [operatorName, setOperatorName] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [formError, setFormError] = useState<string | null>(null);

  const rows = useMemo(() => {
    return banRecords.filter((item) => {
      if (keyword && !item.namespace.toLowerCase().includes(keyword.toLowerCase())) {
        return false;
      }
      return true;
    });
  }, [banRecords, keyword]);

  const handleCreateBan = async () => {
    if (!namespace.trim() || !reason.trim() || !banStartTime.trim() || !operatorName.trim()) {
      setFormError("namespace、原因、开始时间、操作人均为必填。");
      return;
    }

    setSubmitting(true);
    setFormError(null);
    try {
      await createBanRecord({
        namespace: namespace.trim(),
        reason: reason.trim(),
        banStartTime,
        operatorName: operatorName.trim(),
      });
      setOpen(false);
      setNamespace("");
      setReason("");
      setBanStartTime("");
      setOperatorName("");
    } catch (err) {
      setFormError(err instanceof Error ? err.message : "新增封禁记录失败");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="page-container">
      <PageHeader
        kicker="Bans"
        title="封禁记录"
        description="录入和查看封禁记录，新建用弹窗，详情查看用右侧抽屉，不混在一起。"
        actions={<Button variant="primary" onClick={() => setOpen(true)}>新增封禁</Button>}
      />

      <SurfaceCard>
        <div className="toolbar">
          <Field label="namespace">
            <Input placeholder="按 namespace 搜索" value={keyword} onChange={(event) => setKeyword(event.target.value)} />
          </Field>
          <Field label="操作人">
            <Input placeholder="例如：Alice" />
          </Field>
        </div>
      </SurfaceCard>

      <SurfaceCard className="data-table-wrap" padded={false}>
        {rows.length > 0 ? (
          <table className="data-table">
            <thead>
              <tr>
                <th>namespace</th>
                <th>原因</th>
                <th>开始时间</th>
                <th>操作人</th>
                <th>状态</th>
                <th>操作</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((item) => (
                <tr key={item.id}>
                  <td>
                    <button className="namespace-link table-row-button" onClick={() => navigate(`/namespaces/${item.namespace}`)} type="button">
                      {item.namespace}
                    </button>
                  </td>
                  <td>
                    <button className="table-row-button" onClick={() => setSelected(item)} type="button">
                      {item.reason}
                    </button>
                  </td>
                  <td>{item.banStartTime}</td>
                  <td>{item.operatorName}</td>
                  <td>
                    <StatusPill tone="warn">永久封禁</StatusPill>
                  </td>
                  <td>
                    <Button variant="ghost" onClick={() => setSelected(item)}>
                      查看
                    </Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        ) : (
          <div style={{ padding: 20 }}>
            <EmptyState
              title="当前没有封禁记录"
              description="可以直接新增一条封禁记录。"
              action={<Button variant="primary" onClick={() => setOpen(true)}>新增封禁</Button>}
            />
          </div>
        )}
      </SurfaceCard>

      <Drawer
        description="封禁详情在右侧查看，避免把新增表单和详情揉在一个抽屉里。"
        onClose={() => setSelected(null)}
        open={Boolean(selected)}
        title={selected ? selected.namespace : ""}
      >
        {selected ? (
          <>
            <DetailList
              items={[
                { label: "namespace", value: selected.namespace },
                { label: "原因", value: selected.reason },
                { label: "开始时间", value: selected.banStartTime },
                { label: "操作人", value: selected.operatorName },
                { label: "状态", value: "永久封禁" },
              ]}
            />
            <div className="button-row" style={{ marginTop: 20 }}>
              <Button variant="secondary" onClick={() => navigate(`/namespaces/${selected.namespace}`)}>
                查看 namespace 详情
              </Button>
              <Button variant="danger" onClick={() => setPendingDelete(selected)}>
                删除记录
              </Button>
            </div>
          </>
        ) : null}
      </Drawer>

      <Modal
        description="保留简单字段，和当前后端接口一致。"
        onClose={() => {
          setOpen(false);
          setFormError(null);
        }}
        open={open}
        title="新增封禁"
      >
        <div className="panel-stack">
          <Field label="namespace">
            <Input placeholder="例如：prod-finance" value={namespace} onChange={(event) => setNamespace(event.target.value)} />
          </Field>
          <Field label="原因">
            <Input placeholder="简要说明封禁原因" value={reason} onChange={(event) => setReason(event.target.value)} />
          </Field>
          <Field label="开始时间">
            <Input type="datetime-local" value={banStartTime} onChange={(event) => setBanStartTime(event.target.value)} />
          </Field>
          <Field label="操作人">
            <Input placeholder="例如：Alice" value={operatorName} onChange={(event) => setOperatorName(event.target.value)} />
          </Field>
          {formError ? <div className="muted-text" style={{ color: "#b42318" }}>{formError}</div> : null}
          <div className="button-row">
            <Button variant="primary" onClick={() => void handleCreateBan()}>
              {submitting ? "保存中..." : "保存封禁记录"}
            </Button>
            <Button
              variant="secondary"
              onClick={() => {
                setOpen(false);
                setFormError(null);
              }}
            >
              取消
            </Button>
          </div>
        </div>
      </Modal>

      <ConfirmModal
        description={pendingDelete ? `删除后将从当前前端列表中移除 namespace ${pendingDelete.namespace} 的封禁记录。` : ""}
        onClose={() => setPendingDelete(null)}
        onConfirm={() => {
          if (!pendingDelete) return;
          void deleteBanRecord(pendingDelete.namespace).then(() => {
            if (selected?.id === pendingDelete.id) {
              setSelected(null);
            }
            setPendingDelete(null);
          });
        }}
        open={Boolean(pendingDelete)}
        title="删除封禁记录"
      />
    </div>
  );
}
