function divide(filename)
    green=[206 255 100];
    yellow=[255 254 77];
    orange=[255 200 69];
    ocean=[198 255 216];

    A1=imread(filename);

    tresholds = [0.02 0.02 0.06];
    
    for i=1:3
        A(:,:,i)=double(edge(A1(:,:,i), 'sobel', tresholds(i)));
    end

    BW = max(max(A(:,:,1), A(:,:,2)), A(:,:,3)) * 255;
    [sz1,sz2] = size(BW);
    BW(1,:)=255;
    BW(:,sz2)=255;
    
    [B,L,N] = bwboundaries(BW,'holes');
    stat = regionprops(L,'Area','PixelIdxList','BoundingBox');
    
    
    K = zeros(sz1,sz2);
    for k=2:size(stat,1)
        if stat(k).Area > 800
                BWF = zeros(sz1,sz2);
                BWF(stat(k).PixelIdxList) = 255;

                A2 = A1;
                rA=[];
                gA=[];
                bA=[];
                for l=1:sz1
                    for m=1:sz2
                        if BWF(l,m) == 0
                            A2(l,m,1)=0;
                            A2(l,m,2)=0;
                            A2(l,m,3)=0;
                        else
                            rA=[rA;A2(l,m,1)];
                            gA=[gA;A2(l,m,2)];
                            bA=[bA;A2(l,m,3)];
                        end
                    end
                end
                r=median(rA);
                g=median(gA);
                b=median(bA);

                if r > 225 && g > 225 && b > 225
                    continue
                end
                
                gray = rgb2gray(A2);
                bin = imbinarize(gray);
                filled = imfill(bin, 'holes');
                res = medfilt2(filled);
                
                if abs(r-green(1))<10 && abs(g-green(2))<10 && abs(b-green(3))<10
                    imwrite(res,sprintf('images/green/image%d.jpg',k));
                    continue
                end
                if abs(r-orange(1))<10 && abs(g-orange(2))<10 && abs(b-orange(3))<10
                    imwrite(res,sprintf('images/orange/image%d.jpg',k));
                    continue
                end
                if abs(r-yellow(1))<10 && abs(g-yellow(2))<10 && abs(b-yellow(3))<10
                    imwrite(res,sprintf('images/yellow/image%d.jpg',k));
                    continue
                end
                if abs(r-ocean(1))<15 && abs(g-ocean(2))<15 && abs(b-ocean(3))<15
                    imwrite(res,sprintf('images/ocean/image%d.jpg',k));
                    continue
                end
        end
    end
end