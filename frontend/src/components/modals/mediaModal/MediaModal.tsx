import styles from "./MediaModal.module.scss";
import Modal from "../../modal/Modal";
import InputContainer from "../../inputs/inputContainer/InputContainer";

const MediaModal = ({ isOpen, setIsOpen, content, setContent, type }: any) => {
  const onClose = () => {
    setIsOpen(false);
  };

  const onSave = async () => {
    for (const key in content) {
      fetch(`/api/settings/${key}`, {
        method: "PUT",
        headers: {
          "Content-Type": "application/json",
          Authorization: `Bearer ${localStorage.getItem("token")}`,
        },
        body: JSON.stringify({ id: key, value: String(content[key]) }),
      });
    }
    setIsOpen(false);
  };

  let title = "";
  if (type == "posters") title = "Poster Options";
  if (type == "table") title = "Table Options";
  if (type == "overview") title = "Overview Options";
  if (!isOpen) return null;
  return (
    <Modal
      isOpen={isOpen}
      setIsOpen={setIsOpen}
      onSave={onSave}
      title={title}
      onClose={onClose}
    >
      <div className={styles.content}>
        <div className={styles.left}>
          {type === "posters" && (
            <>
              <InputContainer
                label="Poster Size"
                type="select"
                selected={content.mediaPosterPosterSize}
                onChange={(e: any) => {
                  setContent({
                    ...content,
                    mediaPosterPosterSize: e.target.value,
                  });
                }}
              >
                <option value={"small"}>Small</option>
                <option value={"medium"}>Medium</option>
                <option value={"large"}>Large</option>
              </InputContainer>
              <InputContainer
                label="Detailed Progress Bar"
                type="checkbox"
                checked={
                  (content?.mediaPosterDetailedProgressBar == "true")
                }
                helpText="Show text on progress bar"
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaPosterDetailedProgressBar: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Title"
                type="checkbox"
                checked={content?.mediaPosterShowTitle == "true"}
                helpText="Show series title under poster"
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaPosterShowTitle: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Monitored"
                type="checkbox"
                checked={content?.mediaPosterShowMonitored == "true"}
                helpText="Show monitored status under poster"
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaPosterShowMonitored: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Profile"
                type="checkbox"
                checked={content?.mediaPosterShowProfile == "true"}
                helpText="Show codec profile under poster"
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaPosterShowProfile: String(e.target.checked),
                  })
                }
              />
            </>
          )}
          {type === "table" && (
            <>
              <InputContainer
                label="Show Type"
                type="checkbox"
                checked={content?.mediaTableShowType == "true"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaTableShowType: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Network"
                type="checkbox"
                checked={content?.mediaTableShowNetwork == "true"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaTableShowNetwork: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Year"
                type="checkbox"
                checked={content?.mediaTableShowYear == "true"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaTableShowYear: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Profile"
                type="checkbox"
                checked={content?.mediaTableShowProfile == "true"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaTableShowProfile: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Seasons"
                type="checkbox"
                checked={content?.mediaTableShowSeasons == "true"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaTableShowSeasons: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Episodes"
                type="checkbox"
                checked={content?.mediaTableShowEpisodes == "true"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaTableShowEpisodes: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Episode Count"
                type="checkbox"
                checked={content?.mediaTableShowEpisodeCount == "true"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaTableShowEpisodeCount: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Year"
                type="checkbox"
                checked={content?.mediaTableShowYear == "true"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaTableShowYear: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Size On Disk"
                type="checkbox"
                checked={content?.mediaTableShowSizeOnDisk == "true"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaTableShowSizeOnDisk: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Space Saved"
                type="checkbox"
                checked={content?.mediaTableShowSizeSaved == "true"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaTableShowSizeSaved: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Genre"
                type="checkbox"
                checked={content?.mediaTableShowGenre == "true"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaTableShowGenre: String(e.target.checked),
                  })
                }
              />
            </>
          )}
          {type == "overview" && (
            <>
              <InputContainer
                label="Poster Size"
                type="select"
                selected={content?.mediaOverviewPosterSize}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaOverviewPosterSize: e.target.value,
                  })
                }
              >
                <option value={"small"}>Small</option>
                <option value={"medium"}>Medium</option>
                <option value={"large"}>Large</option>
              </InputContainer>
              <InputContainer
                label="Detailed Progress Bar"
                type="checkbox"
                checked={
                  content?.mediaOverviewDetailedProgressBar == "true"
                }
                helpText={"Show text on progress bar"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaOverviewDetailedProgressBar: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Monitored"
                type="checkbox"
                checked={content?.mediaOverviewShowMonitored == "true"}
                helpText={"Show monitored in tags"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaOverviewShowMonitored: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Network"
                type="checkbox"
                checked={content?.mediaOverviewShowNetwork == "true"}
                helpText={"Show network in tags"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaOverviewShowNetwork: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Profile"
                type="checkbox"
                checked={content?.mediaOverviewShowProfile == "true"}
                helpText={"Show profile in tags"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaOverviewShowProfile: String(e.target.checked),
                  })
                }
              />

              <InputContainer
                label="Show Season Count"
                type="checkbox"
                checked={content?.mediaOverviewShowSeasonCount == "true"}
                helpText={"Show season count in tags"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaOverviewShowSeasonCount: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Path"
                type="checkbox"
                checked={content?.mediaOverviewShowPath == "true"}
                helpText={"Show path in tags"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaOverviewShowPath: String(e.target.checked),
                  })
                }
              />
              <InputContainer
                label="Show Size On Disk"
                type="checkbox"
                checked={content?.mediaOverviewShowSizeOnDisk == "true"}
                helpText={"Show size on disk in tags"}
                onChange={(e: any) =>
                  setContent({
                    ...content,
                    mediaOverviewShowSizeOnDisk: String(e.target.checked),
                  })
                }
              />
            </>
          )}
        </div>
      </div>
    </Modal>
  );
};
export default MediaModal;
