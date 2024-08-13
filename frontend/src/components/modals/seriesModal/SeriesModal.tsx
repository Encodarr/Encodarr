import styles from "./SeriesModal.module.scss";
import InputSelect from "../../inputs/inputSelect/InputSelect";
import InputCheckbox from "../../inputs/inputCheckbox/InputCheckbox";
import Modal from "../../modal/Modal";

const SeriesModal = ({
  isOpen,
  setIsOpen,
  content,
  setContent,
  profiles,
}: any) => {
  const onSave = async () => {
    await fetch(`/api/series/${content.name}`, {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
        Authorization: `Bearer ${localStorage.getItem("token")}`,
      },
      body: JSON.stringify(content),
    });
    setIsOpen(false);
  };

  const onClose = () => {
    setIsOpen(false);
  };
  if (!isOpen) return null;
  return (
    <Modal
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      title={"Series Options"}
      onClose={onClose}
      onSave={onSave}
    >
      <div className={styles.content}>
        <div className={styles.inputContainer}>
          <label className={styles.label}>Monitored </label>
          <InputCheckbox
            type="checkbox"
            checked={content?.monitored}
            onChange={(e: any) =>
              setContent({
                ...content,
                monitored: e.target.checked,
              })
            }
          />
        </div>
        <div className={styles.inputContainer}>
          <label className={styles.label}>Profile </label>
          <InputSelect
            selected={content.profile_id.value}
            onChange={(e: any) => {
              setContent({ ...content.profile_id, value: e.target.value });
            }}
          >
            {Object.values(profiles)?.map((profile: any, index: number) => (
              <option value={profile.id} key={index}>
                {profile.name}
              </option>
            ))}
          </InputSelect>
        </div>
      </div>
    </Modal>
  );
};
export default SeriesModal;
