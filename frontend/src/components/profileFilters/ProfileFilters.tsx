import styles from "./ProfileFilters.module.scss";
import InputContainer from "../inputs/inputContainer/InputContainer";

const ProfileFilters = ({ content, setContent }: any) => {
	return (
		<div className={styles.section}>
			<div className={styles.row}>
				<div className={styles.item}>
					<InputContainer
						type="select"
						label="Detelecine"
						selected={content?.detelecine}
						onChange={(e: any) => {
							setContent({ ...content, detelecine: e.target.value });
						}}
					>
						<option value="off">Off</option>
						<option value="default">Default</option>
					</InputContainer>
				</div>
			</div>
			<div className={styles.row}>
				<div className={styles.item}>
					<InputContainer
						type="select"
						label="Interlace"
						selected={content?.interlaceDetection}
						onChange={(e: any) => {
							setContent({
								...content,
								interlaceDetection: e.target.value,
							});
						}}
					>
						<option value="off">Off</option>
						<option value="default">Default</option>
						<option value="less sensitive">Less Sensitive</option>
					</InputContainer>
				</div>
			</div>
			<div className={styles.row}>
				<div className={styles.item}>
					<InputContainer
						type="select"
						label="Deinterlace"
						selected={content?.deinterlace}
						onChange={(e: any) => {
							setContent({ ...content, deinterlace: e.target.value });
						}}
					>
						<option value="off">Off</option>
						<option value="yadif">Yadif</option>
						<option value="decomb">Decomb</option>
						<option value="bwdif">Bwdif</option>
					</InputContainer>
				</div>
				{content?.deinterlace != "off" && (
					<div className={styles.item}>
						<InputContainer
							type="select"
							label="Preset"
							selected={content?.deinterlacePreset}
							onChange={(e: any) => {
								setContent({ ...content, deinterlacePreset: e.target.value });
							}}
						>
							{content?.deinterlace == "yadif" && (
								<>
									<option value="default">Default</option>
									<option value="skip spatial check">Skip Spatial Check</option>
									<option value="Bob">Bob</option>
								</>
							)}
							{content?.deinterlace == "decomb" && (
								<>
									<option value="default">Default</option>
									<option value="bob">Bob</option>
									<option value="eedi2">EEDI2</option>
									<option value="eedl2 bob">EEDI2 Bob</option>
								</>
							)}
							{content?.deinterlace == "bwdif" && (
								<>
									<option value="default">Default</option>
									<option value="bob">Bob</option>
								</>
							)}
						</InputContainer>
					</div>
				)}
			</div>
			<div className={styles.row}>
				<div className={styles.item}>
					<InputContainer
						type="select"
						label="Deblock"
						selected={content?.deblock}
						onChange={(e: any) => {
							setContent({ ...content, deblock: e.target.value });
						}}
					>
						<option value="off">Off</option>
						<option value="ultralight">Ultralight</option>
						<option value="light">Light</option>
						<option value="medium">Medium</option>
						<option value="strong">Strong</option>
						<option value="stronger">Stronger</option>
						<option value="very strong">Very Strong</option>
					</InputContainer>
				</div>
				{content?.deblock != "off" && (
					<div className={styles.item}>
						<InputContainer
							type="select"
							label="Tune"
							selected={content?.deblockTune}
							onChange={(e: any) => {
								setContent({ ...content, deblockTune: e.target.value });
							}}
						>
							<option value="small">Small (4x4)</option>
							<option value="medium">Medium (8x8)</option>
							<option value="large">Large (16x16)</option>
						</InputContainer>
					</div>
				)}
			</div>
			<div className={styles.row}>
				<div className={styles.item}>
					<InputContainer
						type="select"
						label="Denoise"
						selected={content?.denoise}
						onChange={(e: any) => {
							setContent({ ...content, denoise: e.target.value });
						}}
					>
						<option value="off">Off</option>
						<option value="nlmeans">NLMeans</option>
						<option value="hqn3d">HQN3D</option>
					</InputContainer>
				</div>
				{content?.denoise != "off" && (
					<>
						<div className={styles.item}>
							<InputContainer
								type="select"
								label="Preset"
								selected={content?.denoisePreset}
								onChange={(e: any) => {
									setContent({ ...content, denoisePreset: e.target.value });
								}}
							>
								<option value="ultralight">Ultralight</option>
								<option value="light">Light</option>
								<option value="medium">Medium</option>
								<option value="strong">Strong</option>
							</InputContainer>
						</div>
						{content?.denoise == "nlmeans" && (
							<div className={styles.item}>
								<InputContainer
									type="select"
									label="Tune"
									selected={content?.denoiseTune}
									onChange={(e: any) => {
										setContent({ ...content, denoiseTune: e.target.value });
									}}
								>
									<option value="film">Film</option>
									<option value="grain">Grain</option>
									<option value="high motion">High Motion</option>
									<option value="animation">Animation</option>
									<option value="tape">Tape</option>
									<option value="sprite">Sprite</option>
								</InputContainer>
							</div>
						)}
					</>
				)}
			</div>
			<div className={styles.row}>
				<div className={styles.item}>
					<InputContainer
						type="select"
						label="Smooth"
						selected={content?.chromaSmooth}
						onChange={(e: any) => {
							setContent({ ...content, chromaSmooth: e.target.value });
						}}
					>
						<option value="off">Off</option>
						<option value="ultralight">Ultralight</option>
						<option value="light">Light</option>
						<option value="medium">Medium</option>
						<option value="strong">Strong</option>
						<option value="stronger">Stronger</option>
						<option value="very strong">Very Strong</option>
					</InputContainer>
				</div>
				<div className={styles.item}>
					{content?.chromaSmooth != "off" && (
						<InputContainer
							type="select"
							label="Tune"
							selected={content?.chromaSmoothTune}
							onChange={(e: any) => {
								setContent({ ...content, chromaSmoothTune: e.target.value });
							}}
						>
							<option value="off">Off</option>
							<option value="tiny">Tiny</option>
							<option value="small">Small</option>
							<option value="medium">Medium</option>
							<option value="wide">Wide</option>
							<option value="very wide">Very Wide</option>
						</InputContainer>
					)}
				</div>
			</div>
			<div className={styles.row}>
				<div className={styles.item}>
					<InputContainer
						type="select"
						label="Sharpen"
						selected={content?.sharpen}
						onChange={(e: any) => {
							setContent({ ...content, sharpen: e.target.value });
						}}
					>
						<option value="off">Off</option>
						<option value="unsharp">Unsharp</option>
						<option value="lapsharp">Lapsharp</option>
					</InputContainer>
				</div>
				<div className={styles.item}>
					{content?.sharpen != "off" && (
						<InputContainer
							type="select"
							label="Preset"
							selected={content?.sharpenPreset}
							onChange={(e: any) => {
								setContent({ ...content, sharpenPreset: e.target.value });
							}}
						>
							<option value="ultralight">Ultralight</option>
							<option value="light">Light</option>
							<option value="medium">Medium</option>
							<option value="strong">Strong</option>
							<option value="stronger">Stronger</option>
							<option value="very strong">Very Strong</option>
						</InputContainer>
					)}
				</div>
			</div>
			<div className={styles.row}>
				<div className={styles.item}>
					<InputContainer
						type="select"
						label="Colorspace"
						selected={content?.colorspace}
						onChange={(e: any) => {
							setContent({ ...content, colorspace: e.target.value });
						}}
					>
						<option value="off">Off</option>
						<option value="bt.2020">BT.2020</option>
						<option value="bt.709">BT.709</option>
						<option value="bt.601 smpte-c">BT.601 SMPTE-C</option>
						<option value="bt.601 ebu">BT.601 EBU</option>
					</InputContainer>
				</div>
			</div>
			<div className={styles.row}>
				<div className={styles.item}>
					<InputContainer
						type="checkbox"
						label="Color"
						checked={content?.grayscale}
						helpText="Grayscale"
						onChange={(e: any) => {
							setContent({ ...content, grayscale: e.target.checked });
						}}
					/>
				</div>
			</div>
		</div>
	);
};
export default ProfileFilters;
